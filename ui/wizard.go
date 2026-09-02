package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vulan1999/go-gin-template-cli/generator"
	inputtext "github.com/vulan1999/go-gin-template-cli/ui/components/input_text"
	"github.com/vulan1999/go-gin-template-cli/ui/components/spinner"
)

type state int

const (
	stateProjectName state = iota
	stateModuleChoice
	stateGenerating
	stateDone
)

type WizardModel struct {
	state        state
	projectModel inputtext.TextInputModel
	moduleModel  ModuleChoice
	spinnerModel spinner.SpinnerModel
	ProjectName  string
	ModuleName   string
	Cancelled    bool
	Err          error
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
			m.state = stateGenerating

			steps := []spinner.Step{
				{
					Title:     "Creating project directories...",
					DoneTitle: "Created project directories",
					Action: func() error {
						return generator.CreateDirectories(m.ProjectName)
					},
				},
				{
					Title:     fmt.Sprintf("Initializing Go module (%s)...", m.ModuleName),
					DoneTitle: fmt.Sprintf("Initialized Go module (%s)", m.ModuleName),
					Action: func() error {
						return generator.InitGoModule(m.ProjectName, m.ModuleName)
					},
				},
				{
					Title:     "Installing Gin and dependencies...",
					DoneTitle: "Installed Gin and dependencies",
					Action: func() error {
						return generator.InstallDependencies(m.ProjectName)
					},
				},
				{
					Title:     "Creating template files...",
					DoneTitle: "Created template files",
					Action: func() error {
						return generator.CreateTemplateFiles(m.ProjectName, m.ModuleName)
					},
				},
			}

			m.spinnerModel = spinner.CreateModel(
				fmt.Sprintf("🚀 Generating Gin project '%s'...", m.ProjectName),
				steps,
			)
			return m, m.spinnerModel.Init()
		}
	case stateGenerating:
		model, cmd = m.spinnerModel.Update(msg)
		m.spinnerModel = model.(spinner.SpinnerModel)
		if m.spinnerModel.Quitting {
			m.Cancelled = true
			return m, tea.Quit
		}
		if m.spinnerModel.Done {
			m.Err = m.spinnerModel.Err
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
	case stateModuleChoice:
		return m.projectModel.View() + "\n\n" + m.moduleModel.View()
	case stateGenerating, stateDone:
		return m.projectModel.View() + "\n\n" + m.moduleModel.View() + "\n\n" + m.spinnerModel.View()
	}
	return ""
}
