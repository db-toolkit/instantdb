package engines

import (
	"context"
	
	"github.com/db-toolkit/instant-db/src/instantdb/internal/types"
)

// Engine defines the interface for database engines
type Engine interface {
	// Start starts a new database instance
	Start(ctx context.Context, config types.Config) (*types.Instance, error)
	
	// Stop stops a running database instance
	Stop(ctx context.Context, instanceID string) error
	
	// Status returns the status of a database instance
	Status(ctx context.Context, instanceID string) (*types.Status, error)
	
	// GetConnectionURL returns the connection URL for an instance
	GetConnectionURL(instanceID string) (string, error)
	
	// List returns all running instances for this engine
	List() ([]*types.Instance, error)
}
