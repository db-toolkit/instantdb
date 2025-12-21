package commands

import (
	"fmt"
	"os"

	"github.com/db-toolkit/instant-db/src/instantdb/internal/engines"
	"github.com/db-toolkit/instant-db/src/instantdb/internal/utils"
	"github.com/spf13/cobra"
)

var (
	Engine       engines.Engine
	MySQLEngine  engines.Engine
	RedisEngine  engines.Engine
)

// InitEngine initializes the database engines
func InitEngine() {
	homeDir, _ := os.UserHomeDir()
	baseDir := homeDir + "/.instant-db/data"
	
	Engine = engines.NewPostgresEngine(baseDir)
	MySQLEngine = engines.NewMySQLEngine(baseDir)
	RedisEngine = engines.NewRedisEngine(baseDir)
}

// GetEngineForInstance returns the appropriate engine for an instance
func GetEngineForInstance(instanceID string) (engines.Engine, error) {
	instance, err := utils.LoadInstance(instanceID)
	if err != nil {
		return nil, fmt.Errorf("instance not found: %w", err)
	}
	
	switch instance.Engine {
	case "mysql":
		return MySQLEngine, nil
	case "postgres":
		return Engine, nil
	case "redis":
		return RedisEngine, nil
	default:
		return nil, fmt.Errorf("unknown engine: %s", instance.Engine)
	}
}

// GetRootCommand returns the root cobra command with all subcommands
func GetRootCommand(version string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "instant-db",
		Short:   "Instant, isolated database instances for development",
		Long:    `A CLI tool that spins up isolated database instances instantly for development, with zero configuration.`,
		Version: version,
	}

	// Initialize engines
	InitEngine()

	// Add all commands
	rootCmd.AddCommand(StartCmd())
	rootCmd.AddCommand(StopCmd())
	rootCmd.AddCommand(PauseCmd())
	rootCmd.AddCommand(ResumeCmd())
	rootCmd.AddCommand(ListCmd())
	rootCmd.AddCommand(URLCmd())
	rootCmd.AddCommand(StatusCmd())

	return rootCmd
}
