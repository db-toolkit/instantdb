package engines

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/db-toolkit/instant-db/src/instantdb/internal/types"
	"github.com/db-toolkit/instant-db/src/instantdb/internal/utils"
)

type MySQLEngine struct {
	baseDir    string
	daemonPath string
}

func NewMySQLEngine(baseDir string) *MySQLEngine {
	execPath, _ := os.Executable()
	daemonPath := filepath.Join(filepath.Dir(execPath), "mysql-daemon")
	return &MySQLEngine{
		baseDir:    baseDir,
		daemonPath: daemonPath,
	}
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

	cmd := exec.Command(e.daemonPath, strconv.Itoa(config.Port))
	cmd.Stdout = nil
	cmd.Stderr = nil
	
	if err := cmd.Start(); err != nil {
		os.RemoveAll(config.DataDir)
		return nil, fmt.Errorf("failed to start daemon: %w", err)
	}

	time.Sleep(500 * time.Millisecond)

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

	cmd := exec.Command(e.daemonPath, strconv.Itoa(instance.Port))
	cmd.Stdout = nil
	cmd.Stderr = nil
	
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start daemon: %w", err)
	}

	time.Sleep(500 * time.Millisecond)

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

