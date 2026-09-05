package ui

import (
	"fmt"

	charmList "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vulan1999/go-gin-template-cli/generator"
	inputtext "github.com/vulan1999/go-gin-template-cli/ui/components/input_text"
	"github.com/vulan1999/go-gin-template-cli/ui/components/list"
	"github.com/vulan1999/go-gin-template-cli/ui/components/spinner"
)

type state int

const (
	stateProjectName state = iota
	stateModuleChoice
	stateDatabaseChoice
	stateToolkitChoice
	stateGenerating
	stateDone
)

type WizardModel struct {
	state         state
	projectModel  inputtext.TextInputModel
	moduleModel   ModuleChoice
	databaseModel list.ListModel
	toolkitModel  list.ListModel
	spinnerModel  spinner.SpinnerModel
	uiText        string
	ProjectName   string
	ModuleName    string
	Database      string
	Toolkit       string
	Cancelled     bool
	Err           error
}

func (m WizardModel) getSteps() []spinner.Step {
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
			Title:     fmt.Sprintf("Initialize Git Repository for project %s", m.ProjectName),
			DoneTitle: "Git Repository Initialized",
			Action: func() error {
				return generator.InitGitRepository(m.ProjectName)
			},
		},
		{
			Title:     "Installing Gin and dependencies...",
			DoneTitle: "Installed Gin and dependencies",
			Action: func() error {
				return generator.InstallDependencies(m.ProjectName)
			},
		},
	}

	if m.Database != "" && m.Database != "None" {
		steps = append(steps, spinner.Step{
			Title:     fmt.Sprintf("Installing %s driver and %s toolkit...", m.Database, m.Toolkit),
			DoneTitle: fmt.Sprintf("Installed %s driver and %s toolkit", m.Database, m.Toolkit),
			Action: func() error {
				return generator.InstallDatabaseDriverAndToolkit(m.ProjectName, m.Database, m.Toolkit)
			},
		})
	}

	steps = append(steps, spinner.Step{
		Title:     "Creating base template files...",
		DoneTitle: "Created base template files",
		Action: func() error {
			return generator.CreateBaseTemplateFiles(m.ProjectName, m.ModuleName, m.Database, m.Toolkit)
		},
	})

	if m.Database != "" && m.Database != "None" {
		steps = append(steps, spinner.Step{
			Title:     "Creating database template files...",
			DoneTitle: "Created database template files",
			Action: func() error {
				return generator.CreateDatabaseTemplateFiles(m.ProjectName, m.ModuleName, m.Database, m.Toolkit)
			},
		})
	}

	return steps
}

func defaultDatabaseList() list.ListModel {
	items := []charmList.Item{
		list.NewItem("PostgreSQL", "PostgreSQL database setup"),
		list.NewItem("MySQL", "MySQL database setup"),
		list.NewItem("SQLite", "Self-contained serverless SQLite database setup"),
		list.NewItem("None", "Skip database integration"),
	}
	return list.CreateListModel(items, "Choose your database integration:")
}

func defaultToolkitList() list.ListModel {
	items := []charmList.Item{
		list.NewItem("database/sql", "Lightweight -> full control, more manual SQL"),
		list.NewItem("sqlx", "Middle ground -> raw SQL with struct scanning helpers"),
		list.NewItem("GORM", "Less boilerplate, support auto-migrations, easier for CRUD operations"),
	}
	return list.CreateListModel(items, "Choose your database integration toolkit:")
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
			m.uiText = fmt.Sprintf("Project Name: %s\n", m.ProjectName)
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
			m.databaseModel = defaultDatabaseList()
			m.state = stateDatabaseChoice
			m.uiText += fmt.Sprintf("Module Name: %s\n", m.ModuleName)
			return m, m.databaseModel.Init()
		}
	case stateDatabaseChoice:
		model, cmd = m.databaseModel.Update(msg)
		m.databaseModel = model.(list.ListModel)
		if m.databaseModel.Quitting {
			m.Cancelled = true
			return m, tea.Quit
		}
		if m.databaseModel.Done {
			m.Database = m.databaseModel.Value()
			m.uiText += fmt.Sprintf("Database: %s\n", m.Database)
			// if choose None => process to generating phase
			// else => process to toolkit choosing phase
			if m.Database == "None" {
				m.state = stateGenerating

				m.spinnerModel = spinner.CreateModel(
					fmt.Sprintf("🚀 Generating Gin project '%s'...", m.ProjectName),
					m.getSteps(),
				)
				return m, m.spinnerModel.Init()
			} else {
				m.state = stateToolkitChoice
				m.toolkitModel = defaultToolkitList()
				return m, m.toolkitModel.Init()
			}
		}
	case stateToolkitChoice:
		model, cmd = m.toolkitModel.Update(msg)
		m.toolkitModel = model.(list.ListModel)
		if m.toolkitModel.Quitting {
			m.Cancelled = true
			return m, tea.Quit
		}
		if m.toolkitModel.Done {
			m.Toolkit = m.toolkitModel.Value()
			m.uiText += fmt.Sprintf("Toolkit: %s\n", m.Toolkit)
			m.state = stateGenerating

			m.spinnerModel = spinner.CreateModel(
				fmt.Sprintf("🚀 Generating Gin project '%s'...", m.ProjectName),
				m.getSteps(),
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
		return m.uiText + "\n\n" + m.moduleModel.View()
	case stateDatabaseChoice:
		return m.uiText + "\n\n" + m.databaseModel.View()
	case stateToolkitChoice:
		return m.uiText + "\n\n" + m.toolkitModel.View()
	case stateGenerating, stateDone:
		return m.uiText + "\n\n" + m.spinnerModel.View()
	}
	return ""
}
