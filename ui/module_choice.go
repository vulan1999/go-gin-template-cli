package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vulan1999/go-gin-template-cli/ui/components/choice"
	inputtext "github.com/vulan1999/go-gin-template-cli/ui/components/input_text"
)

type moduleChoiceType int

const (
	choiceBarebone moduleChoiceType = iota
	choiceGithub
	choiceGitlab
)

type step int

const (
	stepSelectConvention step = iota
	stepInputRemoteRepo
	stepInputRemoteUser
)

type ModuleChoice struct {
	step              step
	selectedChoice    moduleChoiceType
	moduleChoiceModel choice.ChoicesModel
	modulePathInput   inputtext.TextInputModel
	ProjectName       string
	ModulePath        string
	remoteRepoUrl     string
	remoteRepoInput   inputtext.TextInputModel
	Done              bool
	Cancelled         bool
}

func NewModuleChoice(projectName string) ModuleChoice {
	moduleChoices := []string{
		"Barebone (Local only -> e.g., " + projectName + ")",
		"Github (Remote Github repo -> e.g., github.com/username/" + projectName + ")",
		"Gitlab (Remote Gitlab repo -> e.g., gitlab.com/username/" + projectName + ")",
	}
	return ModuleChoice{
		step:              stepSelectConvention,
		ProjectName:       projectName,
		moduleChoiceModel: choice.CreateChoicesModel(moduleChoices, "Choose your module naming convention:"),
	}
}

func (m ModuleChoice) Init() tea.Cmd {
	return m.moduleChoiceModel.Init()
}

func (m ModuleChoice) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.step {
	case stepSelectConvention:
		newModel, newCmd := m.moduleChoiceModel.Update(msg)
		m.moduleChoiceModel = newModel.(choice.ChoicesModel)
		cmd = newCmd

		if m.moduleChoiceModel.Quitting {
			m.Cancelled = true
			return m, tea.Quit
		}

		if m.moduleChoiceModel.Done {
			m.selectedChoice = moduleChoiceType(m.moduleChoiceModel.Choice)

			switch m.selectedChoice {
			case choiceBarebone:
				// Barebone: Module path is just the project name
				m.ModulePath = m.ProjectName
				m.Done = true
				return m, tea.Quit

			case choiceGithub:
				m.step = stepInputRemoteUser
				m.modulePathInput = inputtext.CreateModel(
					"Enter your GitHub username or organization:",
					"username or organization",
					true,
				)
				return m, m.modulePathInput.Init()

			case choiceGitlab:
				m.step = stepInputRemoteRepo
				m.remoteRepoInput = inputtext.CreateModel(
					"Enter your GitLab remote URL (press Enter for gitlab.com):",
					"gitlab.com",
					false,
				)
				return m, m.remoteRepoInput.Init()
			}
		}

	case stepInputRemoteRepo:
		newModel, newCmd := m.remoteRepoInput.Update(msg)
		m.remoteRepoInput = newModel.(inputtext.TextInputModel)
		cmd = newCmd

		if m.remoteRepoInput.Quitting {
			m.Cancelled = true
			return m, tea.Quit
		}

		if m.remoteRepoInput.Done {
			remoteRepo := strings.TrimSpace(m.remoteRepoInput.Value())
			remoteRepo = strings.TrimPrefix(remoteRepo, "https://")
			remoteRepo = strings.TrimPrefix(remoteRepo, "http://")
			remoteRepo = strings.Trim(remoteRepo, "/")
			if remoteRepo == "" {
				m.remoteRepoUrl = "gitlab.com"
			} else {
				m.remoteRepoUrl = remoteRepo
			}

			m.remoteRepoInput.TextInput.Blur()
			m.step = stepInputRemoteUser
			m.modulePathInput = inputtext.CreateModel(
				"Enter your GitLab username, group or organization:",
				"username or group",
				true,
			)
			return m, m.modulePathInput.Init()
		}

	case stepInputRemoteUser:
		newModel, newCmd := m.modulePathInput.Update(msg)
		m.modulePathInput = newModel.(inputtext.TextInputModel)
		cmd = newCmd

		if m.modulePathInput.Quitting {
			m.Cancelled = true
			return m, tea.Quit
		}

		if m.modulePathInput.Done {
			owner := strings.Trim(strings.TrimSpace(m.modulePathInput.Value()), "/")
			switch m.selectedChoice {
			case choiceGithub:
				m.ModulePath = fmt.Sprintf("github.com/%s/%s", owner, m.ProjectName)
			case choiceGitlab:
				m.ModulePath = fmt.Sprintf("%s/%s/%s", m.remoteRepoUrl, owner, m.ProjectName)
			}
			m.Done = true
			return m, tea.Quit
		}
	}

	return m, cmd
}

func (m ModuleChoice) View() string {
	switch m.step {
	case stepSelectConvention:
		return m.moduleChoiceModel.View()
	case stepInputRemoteRepo:
		return m.moduleChoiceModel.View() + "\n\n" + m.remoteRepoInput.View()
	case stepInputRemoteUser:
		switch m.selectedChoice {
		case choiceGithub:
			return m.moduleChoiceModel.View() + "\n\n" + m.modulePathInput.View()
		case choiceGitlab:
			return m.moduleChoiceModel.View() + "\n\n" + m.remoteRepoInput.View() + "\n\n" + m.modulePathInput.View()
		}
	}
	return ""
}
