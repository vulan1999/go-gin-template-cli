package spinner

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	pendingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
)

// Step represents a generic sequential task with pending/done descriptions and an action function.
type Step struct {
	Title     string
	DoneTitle string
	Action    func() error
}

// StepDoneMsg signals that a step has completed successfully.
type StepDoneMsg struct {
	Index int
}

// StepErrMsg signals that a step has encountered an error.
type StepErrMsg struct {
	Index int
	Err   error
}

// SpinnerModel is a reusable Bubbletea component to display step-by-step progress with an animated spinner.
type SpinnerModel struct {
	Spinner     spinner.Model
	Title       string
	Steps       []Step
	CurrentStep int
	Done        bool
	Quitting    bool
	Err         error
}

// CreateModel initializes a general-purpose SpinnerModel with a title and list of steps.
func CreateModel(title string, steps []Step) SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	return SpinnerModel{
		Spinner: s,
		Title:   title,
		Steps:   steps,
	}
}

func (m SpinnerModel) Init() tea.Cmd {
	return tea.Batch(
		m.Spinner.Tick,
		m.runCurrentStep(),
	)
}

func (m SpinnerModel) runCurrentStep() tea.Cmd {
	if m.CurrentStep >= len(m.Steps) {
		return nil
	}

	index := m.CurrentStep
	action := m.Steps[index].Action

	return func() tea.Msg {
		if action != nil {
			if err := action(); err != nil {
				return StepErrMsg{Index: index, Err: err}
			}
		}
		return StepDoneMsg{Index: index}
	}
}

func (m SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.Quitting = true
			return m, tea.Quit
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd

	case StepDoneMsg:
		m.CurrentStep = msg.Index + 1
		if m.CurrentStep >= len(m.Steps) {
			m.Done = true
			return m, tea.Quit
		}
		return m, m.runCurrentStep()

	case StepErrMsg:
		m.Err = msg.Err
		m.Done = true
		return m, tea.Quit
	}

	return m, nil
}

func (m SpinnerModel) View() string {
	var sb strings.Builder

	if m.Title != "" {
		sb.WriteString(titleStyle.Render(m.Title) + "\n\n")
	}

	for i, step := range m.Steps {
		doneTitle := step.DoneTitle
		if doneTitle == "" {
			doneTitle = step.Title
		}

		if i < m.CurrentStep {
			sb.WriteString(fmt.Sprintf("  %s %s\n", successStyle.Render("✔"), doneTitle))
		} else if i == m.CurrentStep && m.Err != nil {
			sb.WriteString(fmt.Sprintf("  %s %s\n", errorStyle.Render("✖"), step.Title))
		} else if i == m.CurrentStep {
			sb.WriteString(fmt.Sprintf("  %s %s\n", m.Spinner.View(), step.Title))
		} else {
			sb.WriteString(fmt.Sprintf("  %s %s\n", pendingStyle.Render("-"), step.Title))
		}
	}

	if m.Err != nil {
		sb.WriteString("\n" + errorStyle.Render(fmt.Sprintf("Error: %v", m.Err)) + "\n")
	}

	return sb.String()
}
