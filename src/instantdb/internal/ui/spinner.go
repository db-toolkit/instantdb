package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type spinnerModel struct {
	spinner  spinner.Model
	message  string
	done     bool
	err      error
	task     func() error
}

func (m spinnerModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.runTask(),
	)
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case taskDoneMsg:
		m.done = true
		m.err = msg.err
		return m, tea.Quit
	}
	return m, nil
}

func (m spinnerModel) View() string {
	if m.done {
		if m.err != nil {
			return ErrorStyle.Render("✗ " + m.message + " failed\n")
		}
		return SuccessStyle.Render("✓ " + m.message + " complete\n")
	}
	return fmt.Sprintf("%s %s\n", m.spinner.View(), InfoStyle.Render(m.message))
}

func (m spinnerModel) runTask() tea.Cmd {
	return func() tea.Msg {
		err := m.task()
		return taskDoneMsg{err: err}
	}
}

type taskDoneMsg struct {
	err error
}

// ShowSpinner displays a spinner while running a task
func ShowSpinner(message string, task func() error) error {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = InfoStyle

	m := spinnerModel{
		spinner: s,
		message: message,
		task:    task,
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	if finalModel, ok := finalModel.(spinnerModel); ok {
		return finalModel.err
	}

	return nil
}

// ShowSpinnerWithDelay shows spinner with artificial delay for UX
func ShowSpinnerWithDelay(message string, task func() error, minDuration time.Duration) error {
	return ShowSpinner(message, func() error {
		start := time.Now()
		err := task()
		elapsed := time.Since(start)
		
		if elapsed < minDuration {
			time.Sleep(minDuration - elapsed)
		}
		
		return err
	})
}
