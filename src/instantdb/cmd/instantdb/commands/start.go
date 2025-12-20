package commands

import (
	"context"
	"fmt"

	"github.com/db-toolkit/instant-db/src/instantdb/internal/types"
	"github.com/db-toolkit/instant-db/src/instantdb/internal/ui"
	"github.com/spf13/cobra"
)

var (
	startName     string
	startPort     int
	startPersist  bool
	startUsername string
	startPassword string
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
	cmd.Flags().StringVarP(&startUsername, "username", "u", "", "Database username")
	cmd.Flags().StringVar(&startPassword, "password", "", "Database password")

	return cmd
}

func runStart(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Prompt for username if not provided
	if startUsername == "" {
		startUsername = ui.PromptString("Enter database username", "postgres")
	}

	// Prompt for password if not provided
	if startPassword == "" {
		startPassword = ui.PromptPassword("Enter database password", "postgres")
	}

	config := types.Config{
		Name:     startName,
		Port:     startPort,
		Persist:  startPersist,
		Username: startUsername,
		Password: startPassword,
	}

	var instance *types.Instance
	
	// Show spinner while starting
	err := ui.ShowSpinner("Starting PostgreSQL instance", func() error {
		var err error
		instance, err = Engine.Start(ctx, config)
		return err
	})

	if err != nil {
		return fmt.Errorf("failed to start instance: %w", err)
	}

	// Render instance details
	fmt.Println(ui.RenderInstanceDetails(instance))

	return nil
}
