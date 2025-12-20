package commands

import (
	"context"
	"fmt"

	"github.com/db-toolkit/instant-db/src/instantdb/internal/ui"
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

	engine, err := GetEngineForInstance(instanceID)
	if err != nil {
		return err
	}

	// Show spinner while stopping
	err = ui.ShowSpinner(fmt.Sprintf("Stopping instance %s", instanceID), func() error {
		return engine.Stop(ctx, instanceID)
	})

	if err != nil {
		return fmt.Errorf("failed to stop instance: %w", err)
	}

	fmt.Println(ui.SuccessStyle.Render("✅ Instance stopped successfully!\n"))

	return nil
}
