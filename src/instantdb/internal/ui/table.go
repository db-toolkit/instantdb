package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
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

	// Table
	for _, instance := range instances {
		// Instance name with bullet
		name := SuccessStyle.Render("  • " + instance.Name)
		b.WriteString(name + "\n")

		// Details
		b.WriteString(fmt.Sprintf("    %s %s\n", 
			MutedStyle.Render("ID:    "), 
			ValueStyle.Render(instance.ID)))
		b.WriteString(fmt.Sprintf("    %s %s\n", 
			MutedStyle.Render("Port:  "), 
			ValueStyle.Render(fmt.Sprintf("%d", instance.Port))))
		b.WriteString(fmt.Sprintf("    %s %s\n", 
			MutedStyle.Render("Status:"), 
			SuccessStyle.Render(instance.Status)))
		b.WriteString("\n")
	}

	return b.String()
}

// RenderInstanceDetails renders detailed instance information
func RenderInstanceDetails(instance *types.Instance) string {
	var b strings.Builder

	b.WriteString(SuccessStyle.Render("✅ PostgreSQL instance started successfully!\n\n"))

	// Create a box for the details
	details := []string{
		fmt.Sprintf("%s  %s", LabelStyle.Render("Instance ID:"), ValueStyle.Render(instance.ID)),
		fmt.Sprintf("%s  %s", LabelStyle.Render("Name:       "), ValueStyle.Render(instance.Name)),
		fmt.Sprintf("%s  %s", LabelStyle.Render("Port:       "), ValueStyle.Render(fmt.Sprintf("%d", instance.Port))),
		fmt.Sprintf("%s  %s", LabelStyle.Render("Username:   "), ValueStyle.Render(instance.Username)),
		fmt.Sprintf("%s  %s", LabelStyle.Render("Password:   "), ValueStyle.Render(instance.Password)),
		fmt.Sprintf("%s  %s", LabelStyle.Render("Connection: "), 
			ValueStyle.Render(fmt.Sprintf("postgresql://%s:%s@localhost:%d/postgres", 
				instance.Username, instance.Password, instance.Port))),
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor).
		Padding(1, 2)

	b.WriteString(boxStyle.Render(strings.Join(details, "\n")))
	b.WriteString("\n\n")

	// Tips
	b.WriteString(InfoStyle.Render(fmt.Sprintf("💡 Get connection URL: instant-db url %s\n", instance.ID)))
	b.WriteString(InfoStyle.Render(fmt.Sprintf("💡 Stop instance:      instant-db stop %s\n", instance.ID)))

	return b.String()
}

// RenderStatus renders instance status
func RenderStatus(instanceID string, status *types.Status) string {
	var b strings.Builder

	b.WriteString(TitleStyle.Render(fmt.Sprintf("📊 Instance Status: %s\n\n", instanceID)))

	runningIcon := "❌"
	runningText := "No"
	if status.Running {
		runningIcon = "✅"
		runningText = "Yes"
	}

	healthyIcon := "❌"
	healthyText := "No"
	if status.Healthy {
		healthyIcon = "✅"
		healthyText = "Yes"
	}

	b.WriteString(fmt.Sprintf("  %s %s  %s\n", 
		LabelStyle.Render("Running:"), runningIcon, ValueStyle.Render(runningText)))
	b.WriteString(fmt.Sprintf("  %s %s  %s\n", 
		LabelStyle.Render("Healthy:"), healthyIcon, ValueStyle.Render(healthyText)))
	b.WriteString(fmt.Sprintf("  %s  %s\n\n", 
		LabelStyle.Render("Message:"), MutedStyle.Render(status.Message)))

	return b.String()
}
