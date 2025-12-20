package commands

import (
	"context"
	"fmt"

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

	status, err := Engine.Status(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	fmt.Printf("\n📊 Instance Status: %s\n\n", instanceID)

	if status.Running {
		fmt.Println("  Running:  ✅ Yes")
	} else {
		fmt.Println("  Running:  ❌ No")
	}

	if status.Healthy {
		fmt.Println("  Healthy:  ✅ Yes")
	} else {
		fmt.Println("  Healthy:  ❌ No")
	}

	fmt.Printf("  Message:  %s\n\n", status.Message)

	return nil
}
