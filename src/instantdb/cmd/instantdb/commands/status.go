package commands

import (
	"context"
	"fmt"

	"github.com/db-toolkit/instant-db/src/instantdb/internal/ui"
	"github.com/spf13/cobra"
)

// StatusCmd returns the status command
func StatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status [instance-id]",
		Short: "Check status of an instance",
		Long:  `Check the health and status of a running instance.`,
		Args:  cobra.ExactArgs(1),
		RunE:  runStatus,
	}
}

func runStatus(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	instanceID := args[0]

	engine, err := GetEngineForInstance(instanceID)
	if err != nil {
		return err
	}

	status, err := engine.Status(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	fmt.Println(ui.RenderStatus(instanceID, status))

	return nil
}
