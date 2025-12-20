package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/db-toolkit/instant-db/src/instantdb/internal/types"
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
	reader := bufio.NewReader(os.Stdin)

	// Prompt for username if not provided
	if startUsername == "" {
		fmt.Print("Enter database username (default: postgres): ")
		input, _ := reader.ReadString('\n')
		startUsername = strings.TrimSpace(input)
		if startUsername == "" {
			startUsername = "postgres"
		}
	}

	// Prompt for password if not provided
	if startPassword == "" {
		fmt.Print("Enter database password (default: postgres): ")
		input, _ := reader.ReadString('\n')
		startPassword = strings.TrimSpace(input)
		if startPassword == "" {
			startPassword = "postgres"
		}
	}

	config := types.Config{
		Name:     startName,
		Port:     startPort,
		Persist:  startPersist,
		Username: startUsername,
		Password: startPassword,
	}

	fmt.Println("\n🚀 Starting PostgreSQL instance...")

	instance, err := Engine.Start(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to start instance: %w", err)
	}

	fmt.Printf("\n✅ PostgreSQL instance started successfully!\n\n")
	fmt.Printf("  Instance ID:  %s\n", instance.ID)
	fmt.Printf("  Name:         %s\n", instance.Name)
	fmt.Printf("  Port:         %d\n", instance.Port)
	fmt.Printf("  Username:     %s\n", instance.Username)
	fmt.Printf("  Password:     %s\n", instance.Password)
	fmt.Printf("  Connection:   postgresql://%s:%s@localhost:%d/postgres\n\n", instance.Username, instance.Password, instance.Port)
	fmt.Printf("💡 Get connection URL: instant-db url %s\n", instance.ID)
	fmt.Printf("💡 Stop instance:      instant-db stop %s\n", instance.ID)

	return nil
}
