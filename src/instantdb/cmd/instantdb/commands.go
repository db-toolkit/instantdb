package main

import (
	"context"
	"fmt"
	"os"

	"github.com/db-toolkit/instant-db/src/instantdb/internal/engines"
	"github.com/db-toolkit/instant-db/src/instantdb/internal/types"
	"github.com/spf13/cobra"
)

var (
	// Flags
	name    string
	port    int
	persist bool
	
	// Engine
	engine engines.Engine
)

func init() {
	// Initialize PostgreSQL engine with default base directory
	homeDir, _ := os.UserHomeDir()
	baseDir := homeDir + "/.instant-db/data"
	engine = engines.NewPostgresEngine(baseDir)
	
	// Start command flags
	startCmd.Flags().StringVarP(&name, "name", "n", "", "Instance name")
	startCmd.Flags().IntVarP(&port, "port", "p", 0, "Port number (auto-assigned if not specified)")
	startCmd.Flags().BoolVar(&persist, "persist", false, "Keep data after stop")
	
	// Add commands
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(urlCmd)
	rootCmd.AddCommand(statusCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a new PostgreSQL instance",
	Long:  `Start a new isolated PostgreSQL instance with automatic configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		
		config := types.Config{
			Name:    name,
			Port:    port,
			Persist: persist,
		}
		
		fmt.Println("🚀 Starting PostgreSQL instance...")
		
		instance, err := engine.Start(ctx, config)
		if err != nil {
			return fmt.Errorf("failed to start instance: %w", err)
		}
		
		fmt.Printf("\n✅ PostgreSQL instance started successfully!\n\n")
		fmt.Printf("  Instance ID:  %s\n", instance.ID)
		fmt.Printf("  Name:         %s\n", instance.Name)
		fmt.Printf("  Port:         %d\n", instance.Port)
		fmt.Printf("  Connection:   postgresql://localhost:%d/postgres\n\n", instance.Port)
		fmt.Printf("💡 Get connection URL: instant-db url %s\n", instance.ID)
		fmt.Printf("💡 Stop instance:      instant-db stop %s\n", instance.ID)
		
		return nil
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop [instance-id]",
	Short: "Stop a running instance",
	Long:  `Stop a running PostgreSQL instance and clean up resources.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		instanceID := args[0]
		
		fmt.Printf("🛑 Stopping instance %s...\n", instanceID)
		
		if err := engine.Stop(ctx, instanceID); err != nil {
			return fmt.Errorf("failed to stop instance: %w", err)
		}
		
		fmt.Println("✅ Instance stopped successfully!")
		
		return nil
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all running instances",
	Long:  `List all running PostgreSQL instances.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		instances, err := engine.List()
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
	},
}

var urlCmd = &cobra.Command{
	Use:   "url [instance-id]",
	Short: "Get connection URL for an instance",
	Long:  `Get the PostgreSQL connection URL for an instance.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		instanceID := args[0]
		
		url, err := engine.GetConnectionURL(instanceID)
		if err != nil {
			return fmt.Errorf("failed to get connection URL: %w", err)
		}
		
		fmt.Println(url)
		
		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status [instance-id]",
	Short: "Check status of an instance",
	Long:  `Check the health and status of a running instance.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		instanceID := args[0]
		
		status, err := engine.Status(ctx, instanceID)
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
	},
}
