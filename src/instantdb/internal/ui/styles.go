package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	PrimaryColor   = lipgloss.Color("#00D9FF")
	SuccessColor   = lipgloss.Color("#00FF87")
	ErrorColor     = lipgloss.Color("#FF5F87")
	WarningColor   = lipgloss.Color("#FFD700")
	MutedColor     = lipgloss.Color("#6C7086")
	
	// Styles
	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(PrimaryColor).
		MarginBottom(1)
	
	SuccessStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(SuccessColor)
	
	ErrorStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ErrorColor)
	
	InfoStyle = lipgloss.NewStyle().
		Foreground(PrimaryColor)
	
	MutedStyle = lipgloss.NewStyle().
		Foreground(MutedColor)
	
	LabelStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF"))
	
	ValueStyle = lipgloss.NewStyle().
		Foreground(PrimaryColor)
)
