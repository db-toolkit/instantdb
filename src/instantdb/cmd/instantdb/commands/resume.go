package commands

import (
	"context"
	"fmt"

	"github.com/db-toolkit/instant-db/src/instantdb/internal/ui"
	"github.com/spf13/cobra"
)

// ResumeCmd returns the resume command
func ResumeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "resume [instance-id]",
		Short: "Resume a paused instance",
		Long:  `Resume a paused PostgreSQL instance.`,
		Args:  cobra.ExactArgs(1),
		RunE:  runResume,
	}
}

func runResume(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	instanceID := args[0]

	engine, err := GetEngineForInstance(instanceID)
	if err != nil {
		return err
	}

	// Show spinner while resuming
	err = ui.ShowSpinner(fmt.Sprintf("Resuming instance %s", instanceID), func() error {
		return engine.Resume(ctx, instanceID)
	})

	if err != nil {
		return fmt.Errorf("failed to resume instance: %w", err)
	}

	fmt.Println(ui.SuccessStyle.Render("✅ Instance resumed successfully!\n"))
	fmt.Println(ui.InfoStyle.Render(fmt.Sprintf("💡 Get connection URL: instant-db url %s\n", instanceID)))

	return nil
}
