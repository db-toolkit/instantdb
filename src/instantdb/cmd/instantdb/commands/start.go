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
	startEngine   string
)

// StartCmd returns the start command
func StartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start a new database instance",
		Long:  `Start a new isolated database instance with automatic configuration.`,
		RunE:  runStart,
	}

	cmd.Flags().StringVarP(&startName, "name", "n", "", "Instance name")
	cmd.Flags().IntVarP(&startPort, "port", "p", 0, "Port number (auto-assigned if not specified)")
	cmd.Flags().BoolVar(&startPersist, "persist", false, "Keep data after stop")
	cmd.Flags().StringVarP(&startUsername, "username", "u", "", "Database username")
	cmd.Flags().StringVar(&startPassword, "password", "", "Database password")
	cmd.Flags().StringVarP(&startEngine, "engine", "e", "", "Database engine (postgres, mysql, redis)")

	return cmd
}

func runStart(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Prompt for engine if not provided
	if startEngine == "" {
		startEngine = ui.PromptSelect("Select database engine", []string{"postgres", "redis"})
		
		// Prompt for name in interactive mode
		if startName == "" {
			startName = ui.PromptString("Enter instance name (optional)", "")
		}
	}

	// Validate engine
	if startEngine != "postgres" && startEngine != "mysql" && startEngine != "redis" {
		return fmt.Errorf("unsupported engine: %s (supported: postgres, mysql, redis)", startEngine)
	}

	// Set defaults if not provided
	if startUsername == "" {
		if startEngine == "mysql" {
			startUsername = "root"
		} else if startEngine == "redis" {
			startUsername = "default"
		} else {
			startUsername = "postgres"
		}
	}

	if startPassword == "" {
		if startEngine == "mysql" {
			startPassword = "password"
		} else if startEngine == "redis" {
			startPassword = ""
		} else {
			startPassword = "postgres"
		}
	}

	config := types.Config{
		Name:     startName,
		Port:     startPort,
		Persist:  startPersist,
		Username: startUsername,
		Password: startPassword,
		Engine:   startEngine,
	}

	var instance *types.Instance
	
	// Show spinner while starting
	err := ui.ShowSpinner(fmt.Sprintf("Starting %s instance", startEngine), func() error {
		var err error
		if startEngine == "mysql" {
			instance, err = MySQLEngine.Start(ctx, config)
		} else if startEngine == "redis" {
			instance, err = RedisEngine.Start(ctx, config)
		} else {
			instance, err = Engine.Start(ctx, config)
		}
		return err
	})

	if err != nil {
		return fmt.Errorf("failed to start instance: %w", err)
	}

	// Render instance details
	fmt.Println(ui.RenderInstanceDetails(instance))

	return nil
}
