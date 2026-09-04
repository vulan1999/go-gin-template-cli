package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestWizardModelFlow_WithDatabaseAndToolkit(t *testing.T) {
	wizard := NewWizard()

	if wizard.state != stateProjectName {
		t.Fatalf("expected initial state %v, got %v", stateProjectName, wizard.state)
	}

	// 1. Enter project name
	for _, r := range "my-gin-app" {
		m, _ := wizard.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		wizard = m.(WizardModel)
	}
	m, _ := wizard.Update(tea.KeyMsg{Type: tea.KeyEnter})
	wizard = m.(WizardModel)

	if wizard.state != stateModuleChoice {
		t.Fatalf("expected state %v after project name, got %v", stateModuleChoice, wizard.state)
	}
	if wizard.ProjectName != "my-gin-app" {
		t.Errorf("expected ProjectName %q, got %q", "my-gin-app", wizard.ProjectName)
	}

	// 2. Select Barebone module choice (press Enter on first option)
	m, _ = wizard.Update(tea.KeyMsg{Type: tea.KeyEnter})
	wizard = m.(WizardModel)

	if wizard.state != stateDatabaseChoice {
		t.Fatalf("expected state %v after module choice, got %v", stateDatabaseChoice, wizard.state)
	}
	if wizard.ModuleName != "my-gin-app" {
		t.Errorf("expected ModuleName %q, got %q", "my-gin-app", wizard.ModuleName)
	}

	// 3. Select Database choice (press Enter on first option: PostgreSQL)
	m, _ = wizard.Update(tea.KeyMsg{Type: tea.KeyEnter})
	wizard = m.(WizardModel)

	if wizard.state != stateToolkitChoice {
		t.Fatalf("expected state %v after database choice, got %v", stateToolkitChoice, wizard.state)
	}
	if wizard.Database != "PostgreSQL" {
		t.Errorf("expected Database %q, got %q", "PostgreSQL", wizard.Database)
	}

	// 4. Select Toolkit choice (press Enter on first option: database/sql)
	m, _ = wizard.Update(tea.KeyMsg{Type: tea.KeyEnter})
	wizard = m.(WizardModel)

	if wizard.state != stateGenerating {
		t.Fatalf("expected state %v after toolkit choice, got %v", stateGenerating, wizard.state)
	}
	if wizard.Toolkit != "database/sql" {
		t.Errorf("expected Toolkit %q, got %q", "database/sql", wizard.Toolkit)
	}
}

func TestWizardModelFlow_NoneDatabase(t *testing.T) {
	wizard := NewWizard()

	// 1. Enter project name
	for _, r := range "app" {
		m, _ := wizard.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		wizard = m.(WizardModel)
	}
	m, _ := wizard.Update(tea.KeyMsg{Type: tea.KeyEnter})
	wizard = m.(WizardModel)

	// 2. Select module choice
	m, _ = wizard.Update(tea.KeyMsg{Type: tea.KeyEnter})
	wizard = m.(WizardModel)

	if wizard.state != stateDatabaseChoice {
		t.Fatalf("expected state %v, got %v", stateDatabaseChoice, wizard.state)
	}

	// Navigate down to "None" (index 3)
	for i := 0; i < 3; i++ {
		m, _ = wizard.Update(tea.KeyMsg{Type: tea.KeyDown})
		wizard = m.(WizardModel)
	}

	// 3. Select None
	m, _ = wizard.Update(tea.KeyMsg{Type: tea.KeyEnter})
	wizard = m.(WizardModel)

	if wizard.state != stateGenerating {
		t.Fatalf("expected state %v after selecting None, got %v", stateGenerating, wizard.state)
	}
	if wizard.Database != "None" {
		t.Errorf("expected Database %q, got %q", "None", wizard.Database)
	}
	if wizard.Toolkit != "" {
		t.Errorf("expected empty Toolkit for None database, got %q", wizard.Toolkit)
	}
}
