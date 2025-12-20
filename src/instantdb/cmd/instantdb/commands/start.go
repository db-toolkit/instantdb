package commands

import (
	"context"
	"fmt"

	"github.com/db-toolkit/instant-db/src/instantdb/internal/types"
	"github.com/spf13/cobra"
)

var (
	startName    string
	startPort    int
	startPersist bool
)

// StartCmd returns the start command
func StartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start a new PostgreSQL instance",
		Long:  `Start a new isolated PostgreSQL instance with automatic configuration.`,
		RunE:  runStart,
	}

	cmd.Flags().StringVarP(&startName, "name", "n", "", "Instance name")
	cmd.Flags().IntVarP(&startPort, "port", "p", 0, "Port number (auto-assigned if not specified)")
	cmd.Flags().BoolVar(&startPersist, "persist", false, "Keep data after stop")

	return cmd
}

func runStart(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	config := types.Config{
		Name:    startName,
		Port:    startPort,
		Persist: startPersist,
	}

	fmt.Println("🚀 Starting PostgreSQL instance...")

	instance, err := Engine.Start(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to start instance: %w", err)
	}

	fmt.Printf("\n✅ PostgreSQL instance started successfully!\n\n")
	fmt.Printf("  Instance ID:  %s\n", instance.ID)
	fmt.Printf("  Name:         %s\n", instance.Name)
	fmt.Printf("  Port:         %d\n", instance.Port)
	fmt.Printf("  Connection:   postgresql://localhost:%d/postgres\n\n", instance.Port)
	fmt.Printf("💡 Get connection URL: instant-db url %s\n", instance.ID)
	fmt.Printf("💡 Stop instance:      instant-db stop %s\n", instance.ID)

	return nil
}
