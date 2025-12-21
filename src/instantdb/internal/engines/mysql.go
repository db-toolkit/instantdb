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
	"runtime"
	"time"

	"github.com/db-toolkit/instant-db/src/instantdb/internal/types"
	"github.com/db-toolkit/instant-db/src/instantdb/internal/utils"
	"github.com/redis/go-redis/v9"
)

type MySQLEngine struct {
	baseDir   string
	binaryDir string
}

func NewMySQLEngine(baseDir string) *MySQLEngine {
	homeDir, _ := os.UserHomeDir()
	binaryDir := filepath.Join(homeDir, ".instant-db-mysql")
	return &MySQLEngine{
		baseDir:   baseDir,
		binaryDir: binaryDir,
	}
}

func (e *MySQLEngine) downloadMySQL() error {
	version := "8.0.35"
	platform := "darwin-universal"
	
	if runtime.GOOS == "linux" {
		platform = "linux-amd64"
	} else if runtime.GOOS == "windows" {
		platform = "windows-amd64"
	}
	
	ext := ".tar.gz"
	if runtime.GOOS == "windows" {
		ext = ".zip"
	}
	
	url := fmt.Sprintf("https://github.com/db-toolkit/instant-db/releases/download/binaries-v1.0.0/mysql-%s-%s%s", version, platform, ext)
	
	tmpFile := filepath.Join(e.binaryDir, "mysql"+ext)
	
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download mysql: %w", err)
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
	if runtime.GOOS == "windows" {
		return fmt.Errorf("windows not yet supported")
	} else {
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
		header, err := tr.Next()
		if err != nil {
			return err
		}

		mysqlBinary := filepath.Join(e.binaryDir, "mysqld")
		outFile, err := os.Create(mysqlBinary)
		if err != nil {
			return err
		}
		defer outFile.Close()

		if _, err := io.Copy(outFile, tr); err != nil {
			return err
		}

		os.Chmod(mysqlBinary, 0755)
	}

	os.Remove(tmpFile)
	return nil
}

func (e *MySQLEngine) ensureMySQL() (string, error) {
	mysqlBinary := filepath.Join(e.binaryDir, "mysqld")
	
	if _, err := os.Stat(mysqlBinary); err == nil {
		return mysqlBinary, nil
	}

	os.MkdirAll(e.binaryDir, 0755)

	fmt.Println("📦 Downloading MySQL binaries (first time only)...")
	if err := e.downloadMySQL(); err != nil {
		return "", fmt.Errorf("failed to setup mysql: %w", err)
	}

	return mysqlBinary, nil
}

func (e *MySQLEngine) Start(ctx context.Context, config types.Config) (*types.Instance, error) {
	instanceID := utils.GenerateID()
	
	if config.Name == "" {
		config.Name = fmt.Sprintf("mysql-%s", instanceID[:8])
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
		config.Username = "root"
	}
	if config.Password == "" {
		config.Password = "password"
	}

	if err := os.MkdirAll(config.DataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	mysqlBinary, err := e.ensureMySQL()
	if err != nil {
		return nil, err
	}

	// Initialize MySQL data directory
	initCmd := exec.Command(mysqlBinary, "--initialize-insecure", "--datadir="+config.DataDir)
	if err := initCmd.Run(); err != nil {
		os.RemoveAll(config.DataDir)
		return nil, fmt.Errorf("failed to initialize mysql: %w", err)
	}

	// Start MySQL
	cmd := exec.Command(mysqlBinary,
		"--datadir="+config.DataDir,
		"--port="+fmt.Sprintf("%d", config.Port),
		"--bind-address=127.0.0.1",
	)
	
	logFile := filepath.Join(config.DataDir, "mysql.log")
	logFd, _ := os.Create(logFile)
	cmd.Stdout = logFd
	cmd.Stderr = logFd
	
	if err := cmd.Start(); err != nil {
		os.RemoveAll(config.DataDir)
		return nil, fmt.Errorf("failed to start mysql: %w", err)
	}

	time.Sleep(2 * time.Second)

	instance := &types.Instance{
		ID:        instanceID,
		Name:      config.Name,
		Engine:    "mysql",
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
		utils.KillProcessOnPort(config.Port)
		os.RemoveAll(config.DataDir)
		return nil, fmt.Errorf("failed to save instance: %w", err)
	}

	return instance, nil
}

func (e *MySQLEngine) Stop(ctx context.Context, instanceID string) error {
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

func (e *MySQLEngine) Pause(ctx context.Context, instanceID string) error {
	instance, err := utils.LoadInstance(instanceID)
	if err != nil {
		return fmt.Errorf("instance not found: %w", err)
	}

	utils.KillProcessOnPort(instance.Port)

	instance.Status = "paused"
	instance.Paused = true
	return utils.SaveInstance(instance)
}

func (e *MySQLEngine) Resume(ctx context.Context, instanceID string) error {
	instance, err := utils.LoadInstance(instanceID)
	if err != nil {
		return fmt.Errorf("instance not found: %w", err)
	}

	mysqlBinary, err := e.ensureMySQL()
	if err != nil {
		return err
	}

	cmd := exec.Command(mysqlBinary,
		"--datadir="+instance.DataDir,
		"--port="+fmt.Sprintf("%d", instance.Port),
		"--bind-address=127.0.0.1",
	)
	
	logFile := filepath.Join(instance.DataDir, "mysql.log")
	logFd, _ := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	cmd.Stdout = logFd
	cmd.Stderr = logFd

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start mysql: %w", err)
	}

	time.Sleep(2 * time.Second)

	instance.PID = cmd.Process.Pid
	instance.Status = "running"
	instance.Paused = false
	return utils.SaveInstance(instance)
}

func (e *MySQLEngine) GetStatus(ctx context.Context, instanceID string) (string, error) {
	instance, err := utils.LoadInstance(instanceID)
	if err != nil {
		return "", fmt.Errorf("instance not found: %w", err)
	}
	return instance.Status, nil
}

func (e *MySQLEngine) GetConnectionURL(instanceID string) (string, error) {
	instance, err := utils.LoadInstance(instanceID)
	if err != nil {
		return "", fmt.Errorf("instance not found: %w", err)
	}
	return fmt.Sprintf("mysql://%s:%s@127.0.0.1:%d/mysql", instance.Username, instance.Password, instance.Port), nil
}

func (e *MySQLEngine) Status(ctx context.Context, instanceID string) (*types.Status, error) {
	instance, err := utils.LoadInstance(instanceID)
	if err != nil {
		return nil, fmt.Errorf("instance not found: %w", err)
	}
	return &types.Status{
		Running: instance.Status == "running",
		Message: instance.Status,
	}, nil
}

func (e *MySQLEngine) List() ([]*types.Instance, error) {
	all, err := utils.ListInstances()
	if err != nil {
		return nil, err
	}
	var mysql []*types.Instance
	for _, inst := range all {
		if inst.Engine == "mysql" {
			mysql = append(mysql, inst)
		}
	}
	return mysql, nil
}

