package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ListCmd returns the list command
func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all running instances",
		Long:  `List all running PostgreSQL instances.`,
		RunE:  runList,
	}
}

func runList(cmd *cobra.Command, args []string) error {
	instances, err := Engine.List()
	if err != nil {
		return fmt.Errorf("failed to list instances: %w", err)
	}

	if len(instances) == 0 {
		fmt.Println("No running instances found.")
		fmt.Println("\n💡 Start a new instance: instant-db start")
		return nil
	}

	fmt.Printf("\n📋 Running Instances (%d)\n\n", len(instances))

	for _, instance := range instances {
		fmt.Printf("  • %s\n", instance.Name)
		fmt.Printf("    ID:     %s\n", instance.ID)
		fmt.Printf("    Port:   %d\n", instance.Port)
		fmt.Printf("    Status: %s\n", instance.Status)
		fmt.Println()
	}

	return nil
}
