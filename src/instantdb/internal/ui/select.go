package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type selectModel struct {
	choices  []string
	cursor   int
	selected string
	done     bool
}

func (m selectModel) Init() tea.Cmd {
	return nil
}

func (m selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.selected = m.choices[0]
			m.done = true
			return m, tea.Quit
		case "enter":
			m.selected = m.choices[m.cursor]
			m.done = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		}
	}
	return m, nil
}

func (m selectModel) View() string {
	if m.done {
		return ""
	}

	s := LabelStyle.Render("Select database engine:") + "\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = InfoStyle.Render("›")
			choice = SuccessStyle.Render(choice)
		} else {
			choice = MutedStyle.Render(choice)
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\n" + MutedStyle.Render("(↑/↓ to move, enter to select)")

	return s
}

// PromptSelect prompts the user to select from a list of options
func PromptSelect(label string, choices []string) string {
	m := selectModel{
		choices: choices,
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return choices[0]
	}

	if finalModel, ok := finalModel.(selectModel); ok {
		return finalModel.selected
	}

	return choices[0]
}
