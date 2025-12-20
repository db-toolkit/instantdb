package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// StopCmd returns the stop command
func StopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop [instance-id]",
		Short: "Stop a running instance",
		Long:  `Stop a running PostgreSQL instance and clean up resources.`,
		Args:  cobra.ExactArgs(1),
		RunE:  runStop,
	}
}

func runStop(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	instanceID := args[0]

	fmt.Printf("🛑 Stopping instance %s...\n", instanceID)

	if err := Engine.Stop(ctx, instanceID); err != nil {
		return fmt.Errorf("failed to stop instance: %w", err)
	}

	fmt.Println("✅ Instance stopped successfully!")

	return nil
}
