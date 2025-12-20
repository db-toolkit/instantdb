package engines

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/db-toolkit/instant-db/src/instantdb/internal/types"
	"github.com/db-toolkit/instant-db/src/instantdb/internal/utils"
)

// PostgresEngine implements the Engine interface for PostgreSQL
type PostgresEngine struct {
	baseDir   string
	instances map[string]*embeddedpostgres.EmbeddedPostgres
}

// NewPostgresEngine creates a new PostgreSQL engine
func NewPostgresEngine(baseDir string) *PostgresEngine {
	return &PostgresEngine{
		baseDir:   baseDir,
		instances: make(map[string]*embeddedpostgres.EmbeddedPostgres),
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
	if config.Username == "" {
		config.Username = "postgres"
	}
	if config.Password == "" {
		config.Password = "postgres"
	}

	// Create data directory
	if err := os.MkdirAll(config.DataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Create embedded postgres instance
	// Binaries will be downloaded automatically to ~/.embedded-postgres-go/
	postgres := embeddedpostgres.NewDatabase(
		embeddedpostgres.DefaultConfig().
			Port(uint32(config.Port)).
			Username(config.Username).
			Password(config.Password).
			DataPath(config.DataDir).
			RuntimePath(filepath.Join(os.TempDir(), "embedded-pg-runtime")).
			StartTimeout(30 * time.Second),
	)

	// Start PostgreSQL
	if err := postgres.Start(); err != nil {
		os.RemoveAll(config.DataDir)
		return nil, fmt.Errorf("failed to start postgres: %w", err)
	}

	// Store instance reference
	e.instances[instanceID] = postgres

	// Create instance metadata
	instance := &types.Instance{
		ID:        instanceID,
		Name:      config.Name,
		Engine:    "postgres",
		Port:      config.Port,
		DataDir:   config.DataDir,
		PID:       0, // embedded-postgres manages the process
		Status:    "running",
		CreatedAt: time.Now().Unix(),
		Persist:   config.Persist,
		Username:  config.Username,
		Password:  config.Password,
	}

	// Save instance metadata
	if err := utils.SaveInstance(instance); err != nil {
		postgres.Stop()
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

	// Stop the postgres instance if we have a reference
	if postgres, exists := e.instances[instanceID]; exists {
		if err := postgres.Stop(); err != nil {
			return fmt.Errorf("failed to stop server: %w", err)
		}
		delete(e.instances, instanceID)
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

// Pause pauses a running PostgreSQL instance
func (e *PostgresEngine) Pause(ctx context.Context, instanceID string) error {
	instance, err := utils.LoadInstance(instanceID)
	if err != nil {
		return fmt.Errorf("instance not found: %w", err)
	}

	if instance.Paused {
		return fmt.Errorf("instance is already paused")
	}

	// Try to stop using our reference first
	if postgres, exists := e.instances[instanceID]; exists {
		if err := postgres.Stop(); err != nil {
			return fmt.Errorf("failed to pause server: %w", err)
		}
		delete(e.instances, instanceID)
	} else {
		// If no reference, we need to find and kill the process
		// This happens when CLI is restarted
		// For now, we'll create a new instance just to stop it
		postgres := embeddedpostgres.NewDatabase(
			embeddedpostgres.DefaultConfig().
				Port(uint32(instance.Port)).
				DataPath(instance.DataDir).
				RuntimePath(filepath.Join(os.TempDir(), "embedded-pg-runtime")),
		)
		// Stop will kill any process using this port/data dir
		if err := postgres.Stop(); err != nil {
			// Ignore error if already stopped
		}
	}

	// Mark as paused and save
	instance.Paused = true
	instance.Status = "paused"
	if err := utils.SaveInstance(instance); err != nil {
		return fmt.Errorf("failed to save instance: %w", err)
	}

	return nil
}

// Resume resumes a paused PostgreSQL instance
func (e *PostgresEngine) Resume(ctx context.Context, instanceID string) error {
	instance, err := utils.LoadInstance(instanceID)
	if err != nil {
		return fmt.Errorf("instance not found: %w", err)
	}

	if !instance.Paused {
		return fmt.Errorf("instance is not paused")
	}

	// Create embedded postgres instance
	postgres := embeddedpostgres.NewDatabase(
		embeddedpostgres.DefaultConfig().
			Port(uint32(instance.Port)).
			Username(instance.Username).
			Password(instance.Password).
			DataPath(instance.DataDir).
			RuntimePath(filepath.Join(os.TempDir(), "embedded-pg-runtime")).
			StartTimeout(30 * time.Second),
	)

	// Start PostgreSQL
	if err := postgres.Start(); err != nil {
		return fmt.Errorf("failed to resume postgres: %w", err)
	}

	// Store instance reference
	e.instances[instanceID] = postgres

	// Mark as running and save
	instance.Paused = false
	instance.Status = "running"
	if err := utils.SaveInstance(instance); err != nil {
		postgres.Stop()
		return fmt.Errorf("failed to save instance: %w", err)
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

	// Check if we have an active reference
	_, exists := e.instances[instanceID]
	
	// Check if data directory exists
	if _, err := os.Stat(instance.DataDir); os.IsNotExist(err) {
		return &types.Status{
			Running: false,
			Healthy: false,
			Message: "data directory not found",
		}, nil
	}

	return &types.Status{
		Running: exists,
		Healthy: exists,
		Message: "ok",
	}, nil
}

// GetConnectionURL returns the connection URL for an instance
func (e *PostgresEngine) GetConnectionURL(instanceID string) (string, error) {
	instance, err := utils.LoadInstance(instanceID)
	if err != nil {
		return "", fmt.Errorf("instance not found: %w", err)
	}

	return fmt.Sprintf("postgresql://%s:%s@localhost:%d/postgres?sslmode=disable", 
		instance.Username, instance.Password, instance.Port), nil
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
