package commands

import (
	"os"

	"github.com/db-toolkit/instant-db/src/instantdb/internal/engines"
	"github.com/spf13/cobra"
)

var Engine engines.Engine

// InitEngine initializes the database engine
func InitEngine() {
	homeDir, _ := os.UserHomeDir()
	baseDir := homeDir + "/.instant-db/data"
	Engine = engines.NewPostgresEngine(baseDir)
}

// GetRootCommand returns the root cobra command with all subcommands
func GetRootCommand(version string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "instant-db",
		Short:   "Instant, isolated database instances for development",
		Long:    `A CLI tool that spins up isolated database instances instantly for development, with zero configuration.`,
		Version: version,
	}

	// Initialize engine
	InitEngine()

	// Add all commands
	rootCmd.AddCommand(StartCmd())
	rootCmd.AddCommand(StopCmd())
	rootCmd.AddCommand(ListCmd())
	rootCmd.AddCommand(URLCmd())
	rootCmd.AddCommand(StatusCmd())

	return rootCmd
}
