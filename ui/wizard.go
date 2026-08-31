package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	inputtext "github.com/vulan1999/go-gin-template-cli/ui/components/input_text"
)

type state int

const (
	stateProjectName state = iota
	stateModuleChoice
	stateDone
)

type WizardModel struct {
	state        state
	projectModel inputtext.TextInputModel
	moduleModel  ModuleChoice
	ProjectName  string
	ModuleName   string
	Cancelled    bool
}

func NewWizard() WizardModel {
	return WizardModel{
		state:        stateProjectName,
		projectModel: inputtext.CreateModel("Enter Your Gin Project Name", "Your Gin Project Name", true),
	}
}

func (m WizardModel) Init() tea.Cmd {
	return m.projectModel.Init()
}

func (m WizardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var model tea.Model

	switch m.state {
	case stateProjectName:
		model, cmd = m.projectModel.Update(msg)
		m.projectModel = model.(inputtext.TextInputModel)
		if m.projectModel.Quitting {
			m.Cancelled = true
			return m, tea.Quit
		}
		if m.projectModel.Done {
			m.ProjectName = m.projectModel.Value()
			m.projectModel.TextInput.Blur()
			m.moduleModel = NewModuleChoice(m.ProjectName)
			m.state = stateModuleChoice
			return m, m.moduleModel.Init()
		}
	case stateModuleChoice:
		model, cmd = m.moduleModel.Update(msg)
		m.moduleModel = model.(ModuleChoice)
		if m.moduleModel.Cancelled {
			m.Cancelled = true
			return m, tea.Quit
		}
		if m.moduleModel.Done {
			m.ModuleName = m.moduleModel.ModulePath
			m.state = stateDone
			return m, tea.Quit
		}
	}
	return m, cmd
}

func (m WizardModel) View() string {
	switch m.state {
	case stateProjectName:
		return m.projectModel.View()
	case stateModuleChoice, stateDone:
		return m.projectModel.View() + "\n\n" + m.moduleModel.View()
	}
	return ""
}
