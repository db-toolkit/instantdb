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
	// Get all instances from all engines
	postgresInstances, _ := Engine.List()
	mysqlInstances, _ := MySQLEngine.List()
	redisInstances, _ := RedisEngine.List()
	
	// Combine all instances
	allInstances := append(postgresInstances, mysqlInstances...)
	allInstances = append(allInstances, redisInstances...)

	fmt.Println(ui.RenderInstanceTable(allInstances))

	return nil
}
