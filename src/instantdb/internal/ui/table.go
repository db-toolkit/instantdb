package ui

import (
	"fmt"
	"strings"

	"github.com/db-toolkit/instant-db/src/instantdb/internal/types"
)

// RenderInstanceTable renders a table of instances
func RenderInstanceTable(instances []*types.Instance) string {
	if len(instances) == 0 {
		return MutedStyle.Render("No running instances found.\n\n") +
			InfoStyle.Render("💡 Start a new instance: instant-db start\n")
	}

	var b strings.Builder

	// Title
	title := TitleStyle.Render(fmt.Sprintf("📋 Running Instances (%d)", len(instances)))
	b.WriteString(title + "\n\n")

	// Simple list
	for _, instance := range instances {
		b.WriteString(SuccessStyle.Render(fmt.Sprintf("  • %s\n", instance.Name)))
		b.WriteString(fmt.Sprintf("    ID:     %s\n", instance.ID))
		b.WriteString(fmt.Sprintf("    Port:   %d\n", instance.Port))
		b.WriteString(fmt.Sprintf("    Status: %s\n", instance.Status))
		b.WriteString("\n")
	}

	return b.String()
}

// RenderInstanceDetails renders detailed instance information
func RenderInstanceDetails(instance *types.Instance) string {
	var b strings.Builder

	b.WriteString("\n" + SuccessStyle.Render("✅ PostgreSQL instance started successfully!\n\n"))

	// Simple details
	b.WriteString(fmt.Sprintf("  Instance ID:       %s\n", instance.ID))
	b.WriteString(fmt.Sprintf("  Name:              %s\n", instance.Name))
	b.WriteString(fmt.Sprintf("  Port:              %d\n", instance.Port))
	b.WriteString(fmt.Sprintf("  Username:          %s\n", instance.Username))
	b.WriteString(fmt.Sprintf("  Password:          %s\n", instance.Password))
	b.WriteString(fmt.Sprintf("  Connection String: postgresql://%s:%s@localhost:%d/postgres\n\n", 
		instance.Username, instance.Password, instance.Port))

	// Tips
	b.WriteString(InfoStyle.Render(fmt.Sprintf("💡 Stop instance: instant-db stop %s\n", instance.ID)))

	return b.String()
}

// RenderStatus renders instance status
func RenderStatus(instanceID string, status *types.Status) string {
	var b strings.Builder

	b.WriteString("\n" + TitleStyle.Render(fmt.Sprintf("📊 Instance Status: %s\n\n", instanceID)))

	runningIcon := "❌"
	if status.Running {
		runningIcon = "✅"
	}

	healthyIcon := "❌"
	if status.Healthy {
		healthyIcon = "✅"
	}

	b.WriteString(fmt.Sprintf("  Running:  %s\n", runningIcon))
	b.WriteString(fmt.Sprintf("  Healthy:  %s\n", healthyIcon))
	b.WriteString(fmt.Sprintf("  Message:  %s\n\n", status.Message))

	return b.String()
}
