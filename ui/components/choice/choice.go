package choice

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ChoicesModel struct {
	Choices  []string
	Title    string
	Cursor   int
	Choice   int
	Quitting bool
}

var errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

func CreateChoicesModel(choices []string, title string) ChoicesModel {
	return ChoicesModel{
		Choices: choices,
		Title:   title,
		Cursor:  0,
	}
}

func (m ChoicesModel) Init() tea.Cmd {
	return nil
}

func (m ChoicesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// press up or k to move up
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.Choices)-1 {
				m.Cursor++
			}
		case "enter":
			m.Choice = m.Cursor
			m.Quitting = true
			return m, tea.Quit
		// press ctrl + c or esc to quit
		case "esc", "ctrl+c":
			m.Quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ChoicesModel) View() string {
	s := strings.Builder{}
	if m.Title != "" {
		s.WriteString(m.Title + "\n\n")
	}
	for i, choice := range m.Choices {
		if m.Cursor == i {
			s.WriteString("(*) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(choice)
		s.WriteString("\n")
	}
	s.WriteString("\n(press esc to quit)\n")

	return s.String()
}
