package inputtext

import (
	"errors"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

type TextInputModel struct {
	TextInput textinput.Model
	Title     string
	Required  bool
	Error     error
	Quitting  bool
}

func CreateModel(title string, placeholderText string, required bool) TextInputModel {
	ti := textinput.New()
	ti.Placeholder = placeholderText
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 30

	return TextInputModel{
		TextInput: ti,
		Title:     title,
		Required:  required,
	}
}

func (m TextInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m TextInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.Required && strings.TrimSpace(m.TextInput.Value()) == "" {
				m.Error = errors.New("input cannot be empty. Please enter a value")
				return m, nil
			}
			m.Quitting = true
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			m.Quitting = true
			return m, tea.Quit
		default:
			if m.Error != nil {
				m.Error = nil
			}
		}
	case error:
		m.Error = msg
		return m, nil
	}

	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

func (m TextInputModel) View() string {
	var elements []string
	if m.Title != "" {
		elements = append(elements, m.Title+"\n")
	}
	elements = append(elements, m.TextInput.View())

	if m.Error != nil {
		elements = append(elements, errorStyle.Render("\n"+m.Error.Error()))
	}

	elements = append(elements, "\n(press esc to quit)")

	return lipgloss.JoinVertical(
		lipgloss.Top,
		elements...,
	)
}

func (m TextInputModel) Value() string {
	return m.TextInput.Value()
}
