package engines

import (
	"context"
)

// Engine defines the interface for database engines
type Engine interface {
	// Start starts a new database instance
	Start(ctx context.Context, config Config) (*Instance, error)
	
	// Stop stops a running database instance
	Stop(ctx context.Context, instanceID string) error
	
	// Status returns the status of a database instance
	Status(ctx context.Context, instanceID string) (*Status, error)
	
	// GetConnectionURL returns the connection URL for an instance
	GetConnectionURL(instanceID string) (string, error)
	
	// List returns all running instances for this engine
	List() ([]*Instance, error)
}

// Config holds configuration for starting a database instance
type Config struct {
	Name      string
	Port      int
	DataDir   string
	Persist   bool
	WithData  string
}

// Instance represents a running database instance
type Instance struct {
	ID          string
	Name        string
	Engine      string
	Port        int
	DataDir     string
	PID         int
	Status      string
	CreatedAt   int64
	Persist     bool
}

// Status represents the current status of an instance
type Status struct {
	Running   bool
	Healthy   bool
	Message   string
}
