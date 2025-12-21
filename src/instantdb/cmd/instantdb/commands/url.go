package commands

import (
	"fmt"

	"github.com/db-toolkit/instant-db/src/instantdb/internal/utils"
	"github.com/spf13/cobra"
)

// URLCmd returns the url command
func URLCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "url [instance-name-or-id]",
		Short: "Get connection URL for an instance",
		Long:  `Get the PostgreSQL connection URL for an instance.`,
		Args:  cobra.ExactArgs(1),
		RunE:  runURL,
	}
}

func runURL(cmd *cobra.Command, args []string) error {
	instanceID, err := utils.ResolveInstance(args[0])
	if err != nil {
		return err
	}

	engine, err := GetEngineForInstance(instanceID)
	if err != nil {
		return err
	}

	url, err := engine.GetConnectionURL(instanceID)
	if err != nil {
		return fmt.Errorf("failed to get connection URL: %w", err)
	}

	fmt.Println(url)

	return nil
}
