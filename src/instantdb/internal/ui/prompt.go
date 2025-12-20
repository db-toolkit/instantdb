package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type promptModel struct {
	textInput    textinput.Model
	label        string
	defaultValue string
	value        string
	done         bool
}

func (m promptModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m promptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.value = m.textInput.Value()
			if m.value == "" {
				m.value = m.defaultValue
			}
			m.done = true
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			m.done = true
			m.value = m.defaultValue
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m promptModel) View() string {
	if m.done {
		return ""
	}

	prompt := LabelStyle.Render(m.label)
	if m.defaultValue != "" {
		prompt += MutedStyle.Render(fmt.Sprintf(" (default: %s)", m.defaultValue))
	}
	prompt += ": "

	return "\n" + prompt + m.textInput.View()
}

// PromptString prompts the user for a string input
func PromptString(label, defaultValue string) string {
	ti := textinput.New()
	ti.Placeholder = defaultValue
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 30

	m := promptModel{
		textInput:    ti,
		label:        label,
		defaultValue: defaultValue,
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return defaultValue
	}

	if finalModel, ok := finalModel.(promptModel); ok {
		return finalModel.value
	}

	return defaultValue
}

// PromptPassword prompts the user for a password (masked input)
func PromptPassword(label, defaultValue string) string {
	ti := textinput.New()
	ti.Placeholder = strings.Repeat("*", len(defaultValue))
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '•'
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 30

	m := promptModel{
		textInput:    ti,
		label:        label,
		defaultValue: defaultValue,
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return defaultValue
	}

	if finalModel, ok := finalModel.(promptModel); ok {
		return finalModel.value
	}

	return defaultValue
}
