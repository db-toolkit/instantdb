package engines

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/db-toolkit/instant-db/src/instantdb/internal/types"
	"github.com/db-toolkit/instant-db/src/instantdb/internal/utils"
	"github.com/redis/go-redis/v9"
)

type RedisEngine struct {
	baseDir   string
	binaryDir string
}

func NewRedisEngine(baseDir string) *RedisEngine {
	homeDir, _ := os.UserHomeDir()
	binaryDir := filepath.Join(homeDir, ".instant-db-redis")
	return &RedisEngine{
		baseDir:   baseDir,
		binaryDir: binaryDir,
	}
}

func (e *RedisEngine) downloadRedis() error {
	version := "7.2.4"
	url := fmt.Sprintf("https://download.redis.io/releases/redis-%s.tar.gz", version)
	
	tmpFile := filepath.Join(e.binaryDir, "redis.tar.gz")
	
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download redis: %w", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return err
	}

	// Extract
	file, err := os.Open(tmpFile)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	extractDir := filepath.Join(e.binaryDir, fmt.Sprintf("redis-%s", version))
	
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(e.binaryDir, header.Name)
		
		switch header.Typeflag {
		case tar.TypeDir:
			os.MkdirAll(target, 0755)
		case tar.TypeReg:
			os.MkdirAll(filepath.Dir(target), 0755)
			outFile, err := os.Create(target)
			if err != nil {
				return err
			}
			io.Copy(outFile, tr)
			outFile.Close()
			os.Chmod(target, os.FileMode(header.Mode))
		}
	}

	// Compile Redis
	cmd := exec.Command("make", "-j4")
	cmd.Dir = extractDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to compile redis: %w", err)
	}

	// Copy binary
	srcBinary := filepath.Join(extractDir, "src", "redis-server")
	dstBinary := filepath.Join(e.binaryDir, "redis-server")
	
	input, _ := os.ReadFile(srcBinary)
	os.WriteFile(dstBinary, input, 0755)
	
	os.Remove(tmpFile)
	os.RemoveAll(extractDir)

	return nil
}

func (e *RedisEngine) ensureRedis() (string, error) {
	redisBinary := filepath.Join(e.binaryDir, "redis-server")
	
	if _, err := os.Stat(redisBinary); err == nil {
		return redisBinary, nil
	}

	os.MkdirAll(e.binaryDir, 0755)

	// Check if make is available
	if _, err := exec.LookPath("make"); err != nil {
		return "", fmt.Errorf("'make' not found. Redis requires compilation. Install build tools or use system Redis")
	}

	fmt.Println("📦 Downloading Redis binaries (first time only)...")
	if err := e.downloadRedis(); err != nil {
		return "", fmt.Errorf("failed to setup redis: %w", err)
	}

	return redisBinary, nil
}

func (e *RedisEngine) Start(ctx context.Context, config types.Config) (*types.Instance, error) {
	instanceID := utils.GenerateID()
	
	if config.Name == "" {
		config.Name = fmt.Sprintf("redis-%s", instanceID[:8])
	}
	if config.Port == 0 {
		port, err := utils.GetFreePort()
		if err != nil {
			return nil, fmt.Errorf("failed to allocate port: %w", err)
		}
		config.Port = port
	}
	if config.DataDir == "" {
		config.DataDir = filepath.Join(e.baseDir, instanceID)
	}
	if config.Username == "" {
		config.Username = "default"
	}

	if err := os.MkdirAll(config.DataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	redisBinary, err := e.ensureRedis()
	if err != nil {
		return nil, err
	}

	configFile := filepath.Join(config.DataDir, "redis.conf")
	configContent := fmt.Sprintf(`port %d
dir %s
daemonize no
save 900 1
save 300 10
save 60 10000
dbfilename dump.rdb
`, config.Port, config.DataDir)

	if config.Password != "" {
		configContent += fmt.Sprintf("requirepass %s\n", config.Password)
	}

	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write config: %w", err)
	}

	cmd := exec.Command(redisBinary, configFile)
	cmd.Dir = config.DataDir
	
	logFile := filepath.Join(config.DataDir, "redis.log")
	logFd, _ := os.Create(logFile)
	cmd.Stdout = logFd
	cmd.Stderr = logFd

	if err := cmd.Start(); err != nil {
		os.RemoveAll(config.DataDir)
		return nil, fmt.Errorf("failed to start redis: %w", err)
	}

	if err := e.waitForReady(config.Port, config.Password); err != nil {
		cmd.Process.Kill()
		os.RemoveAll(config.DataDir)
		return nil, fmt.Errorf("redis failed to start: %w", err)
	}

	instance := &types.Instance{
		ID:        instanceID,
		Name:      config.Name,
		Engine:    "redis",
		Port:      config.Port,
		DataDir:   config.DataDir,
		PID:       cmd.Process.Pid,
		Status:    "running",
		CreatedAt: time.Now().Unix(),
		Persist:   config.Persist,
		Username:  config.Username,
		Password:  config.Password,
	}

	if err := utils.SaveInstance(instance); err != nil {
		cmd.Process.Kill()
		os.RemoveAll(config.DataDir)
		return nil, fmt.Errorf("failed to save instance: %w", err)
	}

	return instance, nil
}

func (e *RedisEngine) waitForReady(port int, password string) error {
	opts := &redis.Options{
		Addr: fmt.Sprintf("127.0.0.1:%d", port),
	}
	if password != "" {
		opts.Password = password
	}

	client := redis.NewClient(opts)
	defer client.Close()

	for i := 0; i < 30; i++ {
		if err := client.Ping(context.Background()).Err(); err == nil {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	return fmt.Errorf("timeout waiting for redis")
}

func (e *RedisEngine) Stop(ctx context.Context, instanceID string) error {
	instance, err := utils.LoadInstance(instanceID)
	if err != nil {
		return fmt.Errorf("instance not found: %w", err)
	}

	utils.KillProcessOnPort(instance.Port)

	if !instance.Persist {
		os.RemoveAll(instance.DataDir)
	}

	metaFile := filepath.Join(os.Getenv("HOME"), ".instant-db", instanceID+".json")
	os.Remove(metaFile)
	return nil
}

func (e *RedisEngine) Pause(ctx context.Context, instanceID string) error {
	instance, err := utils.LoadInstance(instanceID)
	if err != nil {
		return fmt.Errorf("instance not found: %w", err)
	}

	utils.KillProcessOnPort(instance.Port)

	instance.Status = "paused"
	instance.Paused = true
	return utils.SaveInstance(instance)
}

func (e *RedisEngine) Resume(ctx context.Context, instanceID string) error {
	instance, err := utils.LoadInstance(instanceID)
	if err != nil {
		return fmt.Errorf("instance not found: %w", err)
	}

	redisBinary, err := e.ensureRedis()
	if err != nil {
		return err
	}

	configFile := filepath.Join(instance.DataDir, "redis.conf")
	cmd := exec.Command(redisBinary, configFile)
	cmd.Dir = instance.DataDir
	
	logFile := filepath.Join(instance.DataDir, "redis.log")
	logFd, _ := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	cmd.Stdout = logFd
	cmd.Stderr = logFd

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start redis: %w", err)
	}

	if err := e.waitForReady(instance.Port, instance.Password); err != nil {
		cmd.Process.Kill()
		return fmt.Errorf("redis failed to start: %w", err)
	}

	instance.PID = cmd.Process.Pid
	instance.Status = "running"
	instance.Paused = false
	return utils.SaveInstance(instance)
}

func (e *RedisEngine) GetConnectionURL(instanceID string) (string, error) {
	instance, err := utils.LoadInstance(instanceID)
	if err != nil {
		return "", fmt.Errorf("instance not found: %w", err)
	}
	
	if instance.Password != "" {
		return fmt.Sprintf("redis://:%s@127.0.0.1:%d", instance.Password, instance.Port), nil
	}
	return fmt.Sprintf("redis://127.0.0.1:%d", instance.Port), nil
}

func (e *RedisEngine) Status(ctx context.Context, instanceID string) (*types.Status, error) {
	instance, err := utils.LoadInstance(instanceID)
	if err != nil {
		return nil, fmt.Errorf("instance not found: %w", err)
	}
	return &types.Status{
		Running: instance.Status == "running",
		Message: instance.Status,
	}, nil
}

func (e *RedisEngine) List() ([]*types.Instance, error) {
	all, err := utils.ListInstances()
	if err != nil {
		return nil, err
	}
	var redis []*types.Instance
	for _, inst := range all {
		if inst.Engine == "redis" {
			redis = append(redis, inst)
		}
	}
	return redis, nil
}
