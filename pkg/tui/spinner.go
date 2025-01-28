package tui

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	spinner    spinner.Model
	program    *tea.Program
	text       string
	loading    bool
	cancelFunc context.CancelFunc
}

func (m Model) Init() tea.Cmd {
	// Start the spinner
	return spinner.Tick
}

func (m *Model) Stop() {
	m.program.ReleaseTerminal()
	m.program.Quit()
	m.cancelFunc()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" {
			// Exit on 'q'
			return m, tea.Quit
		}
	case spinner.TickMsg:
		// Update the spinner
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Model) View() string {
	if m.loading {
		return fmt.Sprintf("%s %s (Press 'q' to quit)", m.text, m.spinner.View())
	}
	return "Done! Press 'q' to exit."
}

func Spinner(ctx context.Context, cancel context.CancelFunc, text string) *Model {
	// Create a spinner Model
	s := spinner.New()
	s.Spinner = spinner.Moon

	// Initialize the program Model
	m := Model{
		spinner:    s,
		text:       text,
		loading:    true,
		cancelFunc: cancel,
	}

	// Run the program
	p := tea.NewProgram(m)
	m.program = p
	go func() {
		p.Run()
	}()

	return &m
}
