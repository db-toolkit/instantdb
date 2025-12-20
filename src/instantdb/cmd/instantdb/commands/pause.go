package commands

import (
	"context"
	"fmt"

	"github.com/db-toolkit/instant-db/src/instantdb/internal/ui"
	"github.com/spf13/cobra"
)

// PauseCmd returns the pause command
func PauseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pause [instance-id]",
		Short: "Pause a running instance",
		Long:  `Pause a running PostgreSQL instance. The data is preserved and can be resumed later.`,
		Args:  cobra.ExactArgs(1),
		RunE:  runPause,
	}
}

func runPause(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	instanceID := args[0]

	engine, err := GetEngineForInstance(instanceID)
	if err != nil {
		return err
	}

	// Show spinner while pausing
	err = ui.ShowSpinner(fmt.Sprintf("Pausing instance %s", instanceID), func() error {
		return engine.Pause(ctx, instanceID)
	})

	if err != nil {
		return fmt.Errorf("failed to pause instance: %w", err)
	}

	fmt.Println(ui.SuccessStyle.Render("✅ Instance paused successfully!\n"))
	fmt.Println(ui.InfoStyle.Render(fmt.Sprintf("💡 Resume instance: instant-db resume %s\n", instanceID)))

	return nil
}
