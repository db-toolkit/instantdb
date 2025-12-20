package engines

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/db-toolkit/instant-db/src/instantdb/internal/types"
	"github.com/db-toolkit/instant-db/src/instantdb/internal/utils"
)

// MySQLEngine implements the Engine interface for MySQL
type MySQLEngine struct {
	baseDir string
}

// NewMySQLEngine creates a new MySQL engine
func NewMySQLEngine(baseDir string) *MySQLEngine {
	return &MySQLEngine{
		baseDir: baseDir,
	}
}

// Start starts a new MySQL instance
func (e *MySQLEngine) Start(ctx context.Context, config types.Config) (*types.Instance, error) {
	// Generate instance ID
	instanceID := utils.GenerateID()
	
	// Set defaults
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

	// Create data directory
	if err := os.MkdirAll(config.DataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Initialize MySQL data directory
	if err := e.initDB(ctx, config.DataDir); err != nil {
		os.RemoveAll(config.DataDir)
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Start MySQL server
	pid, err := e.startServer(ctx, config.DataDir, config.Port, config.Password)
	if err != nil {
		os.RemoveAll(config.DataDir)
		return nil, fmt.Errorf("failed to start server: %w", err)
	}

	// Wait for server to be ready
	if err := e.waitForReady(ctx, config.Port); err != nil {
		utils.KillProcessOnPort(config.Port)
		os.RemoveAll(config.DataDir)
		return nil, fmt.Errorf("server failed to become ready: %w", err)
	}

	// Create instance
	instance := &types.Instance{
		ID:        instanceID,
		Name:      config.Name,
		Engine:    "mysql",
		Port:      config.Port,
		DataDir:   config.DataDir,
		PID:       pid,
		Status:    "running",
		CreatedAt: time.Now().Unix(),
		Persist:   config.Persist,
		Username:  config.Username,
		Password:  config.Password,
	}

	// Save instance metadata
	if err := utils.SaveInstance(instance); err != nil {
		utils.KillProcessOnPort(config.Port)
		os.RemoveAll(config.DataDir)
		return nil, fmt.Errorf("failed to save instance: %w", err)
	}

	return instance, nil
}

// Stop stops a running MySQL instance
func (e *MySQLEngine) Stop(ctx context.Context, instanceID string) error {
	instance, err := utils.LoadInstance(instanceID)
	if err != nil {
		return fmt.Errorf("instance not found: %w", err)
	}

	// Kill the MySQL process
	if err := utils.KillProcessOnPort(instance.Port); err != nil {
		return fmt.Errorf("failed to stop server: %w", err)
	}

	// Clean up data directory if not persistent
	if !instance.Persist {
		if err := os.RemoveAll(instance.DataDir); err != nil {
			return fmt.Errorf("failed to remove data directory: %w", err)
		}
	}

	// Remove instance metadata
	if err := utils.RemoveInstance(instanceID); err != nil {
		return fmt.Errorf("failed to remove instance metadata: %w", err)
	}

	return nil
}

// Pause pauses a running MySQL instance
func (e *MySQLEngine) Pause(ctx context.Context, instanceID string) error {
	instance, err := utils.LoadInstance(instanceID)
	if err != nil {
		return fmt.Errorf("instance not found: %w", err)
	}

	if instance.Paused {
		return fmt.Errorf("instance is already paused")
	}

	// Kill the MySQL process
	if err := utils.KillProcessOnPort(instance.Port); err != nil {
		return fmt.Errorf("failed to stop server: %w", err)
	}

	// Mark as paused and save
	instance.Paused = true
	instance.Status = "paused"
	if err := utils.SaveInstance(instance); err != nil {
		return fmt.Errorf("failed to save instance: %w", err)
	}

	return nil
}

// Resume resumes a paused MySQL instance
func (e *MySQLEngine) Resume(ctx context.Context, instanceID string) error {
	instance, err := utils.LoadInstance(instanceID)
	if err != nil {
		return fmt.Errorf("instance not found: %w", err)
	}

	if !instance.Paused {
		return fmt.Errorf("instance is not paused")
	}

	// Start MySQL server
	pid, err := e.startServer(ctx, instance.DataDir, instance.Port, instance.Password)
	if err != nil {
		return fmt.Errorf("failed to resume server: %w", err)
	}

	// Wait for server to be ready
	if err := e.waitForReady(ctx, instance.Port); err != nil {
		utils.KillProcessOnPort(instance.Port)
		return fmt.Errorf("server failed to become ready: %w", err)
	}

	// Mark as running and save
	instance.Paused = false
	instance.Status = "running"
	instance.PID = pid
	if err := utils.SaveInstance(instance); err != nil {
		utils.KillProcessOnPort(instance.Port)
		return fmt.Errorf("failed to save instance: %w", err)
	}

	return nil
}

// Status returns the status of a MySQL instance
func (e *MySQLEngine) Status(ctx context.Context, instanceID string) (*types.Status, error) {
	instance, err := utils.LoadInstance(instanceID)
	if err != nil {
		return &types.Status{
			Running: false,
			Healthy: false,
			Message: "instance not found",
		}, nil
	}

	// Check if data directory exists
	if _, err := os.Stat(instance.DataDir); os.IsNotExist(err) {
		return &types.Status{
			Running: false,
			Healthy: false,
			Message: "data directory not found",
		}, nil
	}

	// Check if process is running
	running := utils.IsProcessRunning(instance.PID)
	
	return &types.Status{
		Running: running,
		Healthy: running,
		Message: "ok",
	}, nil
}

// GetConnectionURL returns the connection URL for an instance
func (e *MySQLEngine) GetConnectionURL(instanceID string) (string, error) {
	instance, err := utils.LoadInstance(instanceID)
	if err != nil {
		return "", fmt.Errorf("instance not found: %w", err)
	}

	return fmt.Sprintf("mysql://%s:%s@localhost:%d/mysql", 
		instance.Username, instance.Password, instance.Port), nil
}

// List returns all running MySQL instances
func (e *MySQLEngine) List() ([]*types.Instance, error) {
	instances, err := utils.ListInstances()
	if err != nil {
		return nil, err
	}

	// Filter for mysql instances
	var mysqlInstances []*types.Instance
	for _, instance := range instances {
		if instance.Engine == "mysql" {
			mysqlInstances = append(mysqlInstances, instance)
		}
	}

	return mysqlInstances, nil
}

// initDB initializes a MySQL data directory
func (e *MySQLEngine) initDB(ctx context.Context, dataDir string) error {
	cmd := exec.CommandContext(ctx, "mysqld", "--initialize-insecure", "--datadir="+dataDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("mysqld --initialize failed: %w", err)
	}
	
	return nil
}

// startServer starts the MySQL server
func (e *MySQLEngine) startServer(ctx context.Context, dataDir string, port int, password string) (int, error) {
	cmd := exec.CommandContext(
		ctx,
		"mysqld",
		"--datadir="+dataDir,
		"--port="+fmt.Sprintf("%d", port),
		"--socket="+filepath.Join(dataDir, "mysql.sock"),
		"--pid-file="+filepath.Join(dataDir, "mysql.pid"),
	)
	
	// Redirect output to log file
	logFile := filepath.Join(dataDir, "mysql.log")
	f, err := os.Create(logFile)
	if err != nil {
		return 0, fmt.Errorf("failed to create log file: %w", err)
	}
	defer f.Close()
	
	cmd.Stdout = f
	cmd.Stderr = f
	
	if err := cmd.Start(); err != nil {
		return 0, fmt.Errorf("failed to start mysqld: %w", err)
	}
	
	return cmd.Process.Pid, nil
}

// waitForReady waits for MySQL to be ready to accept connections
func (e *MySQLEngine) waitForReady(ctx context.Context, port int) error {
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	
	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for server to be ready")
		case <-ticker.C:
			cmd := exec.CommandContext(ctx, "mysqladmin", "ping", "-P", fmt.Sprintf("%d", port), "-h", "127.0.0.1")
			if cmd.Run() == nil {
				return nil
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
