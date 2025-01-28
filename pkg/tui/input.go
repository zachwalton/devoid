package tui

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func Input(prompt string) string {
	m := initialModel(prompt)
	p := tea.NewProgram(&m)
	p.Run()
	return m.textInput.Value()
}

type (
	errMsg error
)

type inputModel struct {
	prompt    string
	textInput textinput.Model
	response  string
	err       error
}

func initialModel(prompt string) inputModel {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 240
	ti.Width = 0 // unlimited

	return inputModel{
		prompt:    prompt,
		textInput: ti,
		err:       nil,
	}
}

func (m inputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.response = msg.String()
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m inputModel) View() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		m.prompt,
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}
