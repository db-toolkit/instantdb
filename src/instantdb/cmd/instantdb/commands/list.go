package commands

import (
	"fmt"

	"github.com/db-toolkit/instant-db/src/instantdb/internal/ui"
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

	fmt.Println(ui.RenderInstanceTable(instances))

	return nil
}
