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

// PostgresEngine implements the Engine interface for PostgreSQL
type PostgresEngine struct {
	baseDir string
}

// NewPostgresEngine creates a new PostgreSQL engine
func NewPostgresEngine(baseDir string) *PostgresEngine {
	return &PostgresEngine{
		baseDir: baseDir,
	}
}

// Start starts a new PostgreSQL instance
func (e *PostgresEngine) Start(ctx context.Context, config types.Config) (*types.Instance, error) {
	// Generate instance ID
	instanceID := utils.GenerateID()
	
	// Set defaults
	if config.Name == "" {
		config.Name = fmt.Sprintf("postgres-%s", instanceID[:8])
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

	// Create data directory
	if err := os.MkdirAll(config.DataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Initialize PostgreSQL data directory
	if err := e.initDB(ctx, config.DataDir); err != nil {
		os.RemoveAll(config.DataDir)
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Start PostgreSQL server
	pid, err := e.startServer(ctx, config.DataDir, config.Port)
	if err != nil {
		os.RemoveAll(config.DataDir)
		return nil, fmt.Errorf("failed to start server: %w", err)
	}

	// Wait for server to be ready
	if err := e.waitForReady(ctx, config.Port); err != nil {
		e.stopServer(pid)
		os.RemoveAll(config.DataDir)
		return nil, fmt.Errorf("server failed to become ready: %w", err)
	}

	// Create instance
	instance := &types.Instance{
		ID:        instanceID,
		Name:      config.Name,
		Engine:    "postgres",
		Port:      config.Port,
		DataDir:   config.DataDir,
		PID:       pid,
		Status:    "running",
		CreatedAt: time.Now().Unix(),
		Persist:   config.Persist,
	}

	// Save instance metadata
	if err := utils.SaveInstance(instance); err != nil {
		e.stopServer(pid)
		os.RemoveAll(config.DataDir)
		return nil, fmt.Errorf("failed to save instance: %w", err)
	}

	return instance, nil
}

// Stop stops a running PostgreSQL instance
func (e *PostgresEngine) Stop(ctx context.Context, instanceID string) error {
	instance, err := utils.LoadInstance(instanceID)
	if err != nil {
		return fmt.Errorf("instance not found: %w", err)
	}

	// Stop the server
	if err := e.stopServer(instance.PID); err != nil {
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

// Status returns the status of a PostgreSQL instance
func (e *PostgresEngine) Status(ctx context.Context, instanceID string) (*types.Status, error) {
	instance, err := utils.LoadInstance(instanceID)
	if err != nil {
		return &types.Status{
			Running: false,
			Healthy: false,
			Message: "instance not found",
		}, nil
	}

	// Check if process is running
	running := utils.IsProcessRunning(instance.PID)
	if !running {
		return &types.Status{
			Running: false,
			Healthy: false,
			Message: "process not running",
		}, nil
	}

	// Check if server is healthy
	healthy := e.isHealthy(ctx, instance.Port)
	
	return &types.Status{
		Running: running,
		Healthy: healthy,
		Message: "ok",
	}, nil
}

// GetConnectionURL returns the connection URL for an instance
func (e *PostgresEngine) GetConnectionURL(instanceID string) (string, error) {
	instance, err := utils.LoadInstance(instanceID)
	if err != nil {
		return "", fmt.Errorf("instance not found: %w", err)
	}

	return fmt.Sprintf("postgresql://localhost:%d/postgres", instance.Port), nil
}

// List returns all running PostgreSQL instances
func (e *PostgresEngine) List() ([]*types.Instance, error) {
	instances, err := utils.ListInstances()
	if err != nil {
		return nil, err
	}

	// Filter for postgres instances
	var postgresInstances []*types.Instance
	for _, instance := range instances {
		if instance.Engine == "postgres" {
			postgresInstances = append(postgresInstances, instance)
		}
	}

	return postgresInstances, nil
}

// initDB initializes a PostgreSQL data directory
func (e *PostgresEngine) initDB(ctx context.Context, dataDir string) error {
	cmd := exec.CommandContext(ctx, "initdb", "-D", dataDir, "--no-locale", "--encoding=UTF8")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("initdb failed: %w", err)
	}
	
	return nil
}

// startServer starts the PostgreSQL server
func (e *PostgresEngine) startServer(ctx context.Context, dataDir string, port int) (int, error) {
	cmd := exec.CommandContext(
		ctx,
		"postgres",
		"-D", dataDir,
		"-p", fmt.Sprintf("%d", port),
		"-k", dataDir, // Unix socket directory
	)
	
	// Redirect output to log file
	logFile := filepath.Join(dataDir, "postgres.log")
	f, err := os.Create(logFile)
	if err != nil {
		return 0, fmt.Errorf("failed to create log file: %w", err)
	}
	defer f.Close()
	
	cmd.Stdout = f
	cmd.Stderr = f
	
	if err := cmd.Start(); err != nil {
		return 0, fmt.Errorf("failed to start postgres: %w", err)
	}
	
	return cmd.Process.Pid, nil
}

// stopServer stops the PostgreSQL server
func (e *PostgresEngine) stopServer(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process: %w", err)
	}
	
	// Send SIGTERM for graceful shutdown
	if err := process.Signal(os.Interrupt); err != nil {
		return fmt.Errorf("failed to stop process: %w", err)
	}
	
	// Wait for process to exit (with timeout)
	done := make(chan error, 1)
	go func() {
		_, err := process.Wait()
		done <- err
	}()
	
	select {
	case <-time.After(10 * time.Second):
		// Force kill if not stopped
		process.Kill()
		return fmt.Errorf("process did not stop gracefully, killed")
	case err := <-done:
		return err
	}
}

// waitForReady waits for PostgreSQL to be ready to accept connections
func (e *PostgresEngine) waitForReady(ctx context.Context, port int) error {
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	
	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for server to be ready")
		case <-ticker.C:
			if e.isHealthy(ctx, port) {
				return nil
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// isHealthy checks if PostgreSQL is healthy
func (e *PostgresEngine) isHealthy(ctx context.Context, port int) bool {
	cmd := exec.CommandContext(ctx, "pg_isready", "-p", fmt.Sprintf("%d", port))
	return cmd.Run() == nil
}
