package list

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func TestItem(t *testing.T) {
	item := NewItem("PostgreSQL", "Powerful, open-source object-relational database system")
	if item.Title() != "PostgreSQL" {
		t.Errorf("expected Title %q, got %q", "PostgreSQL", item.Title())
	}
	if item.Description() != "Powerful, open-source object-relational database system" {
		t.Errorf("expected Description %q, got %q", "Powerful, open-source object-relational database system", item.Description())
	}
	if item.FilterValue() != "PostgreSQL" {
		t.Errorf("expected FilterValue %q, got %q", "PostgreSQL", item.FilterValue())
	}
}

func TestCreateListModel(t *testing.T) {
	items := []list.Item{
		NewItem("PostgreSQL", "Postgres driver + GORM"),
		NewItem("MySQL", "MySQL driver + GORM"),
		NewItem("SQLite", "SQLite driver + GORM"),
		NewItem("None", "No database configuration"),
	}
	title := "Select Database Integration"
	m := CreateListModel(items, title)

	if m.List.Title != title {
		t.Errorf("expected Title %q, got %q", title, m.List.Title)
	}
	if len(m.List.Items()) != len(items) {
		t.Errorf("expected %d items, got %d", len(items), len(m.List.Items()))
	}
	if m.Done {
		t.Errorf("expected Done to be false initially")
	}
	if m.Quitting {
		t.Errorf("expected Quitting to be false initially")
	}
	if m.SelectedItem != nil {
		t.Errorf("expected SelectedItem to be nil initially")
	}
	if m.Value() != "" {
		t.Errorf("expected Value() to be empty initially, got %q", m.Value())
	}
}

func TestInit(t *testing.T) {
	m := CreateListModel([]list.Item{NewItem("A", "Desc")}, "Title")
	cmd := m.Init()
	if cmd != nil {
		t.Errorf("expected Init to return nil, got %v", cmd)
	}
}

func TestUpdate_Quit(t *testing.T) {
	tests := []struct {
		name   string
		keyMsg tea.KeyMsg
	}{
		{"Ctrl+C", tea.KeyMsg{Type: tea.KeyCtrlC}},
		{"Esc", tea.KeyMsg{Type: tea.KeyEsc}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := CreateListModel([]list.Item{NewItem("A", "Desc")}, "Title")
			newModel, cmd := m.Update(tt.keyMsg)
			lm, ok := newModel.(ListModel)
			if !ok {
				t.Fatalf("expected ListModel, got %T", newModel)
			}
			if !lm.Quitting {
				t.Errorf("expected Quitting to be true")
			}
			if cmd == nil {
				t.Fatalf("expected tea.Quit cmd, got nil")
			}
			msg := cmd()
			if _, isQuit := msg.(tea.QuitMsg); !isQuit {
				t.Errorf("expected QuitMsg, got %T", msg)
			}
		})
	}
}

func TestUpdate_Selection(t *testing.T) {
	items := []list.Item{
		NewItem("PostgreSQL", "Postgres driver + GORM"),
		NewItem("MySQL", "MySQL driver + GORM"),
		NewItem("SQLite", "SQLite driver + GORM"),
	}
	m := CreateListModel(items, "Select Database")

	// Select current (first) item with enter
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	lm := newModel.(ListModel)

	if !lm.Done {
		t.Errorf("expected Done to be true after Enter")
	}
	if lm.SelectedItem == nil {
		t.Fatalf("expected SelectedItem to not be nil")
	}
	if lm.Value() != "PostgreSQL" {
		t.Errorf("expected Value() to be %q, got %q", "PostgreSQL", lm.Value())
	}
}

func TestUpdate_WindowSize(t *testing.T) {
	m := CreateListModel([]list.Item{NewItem("A", "Desc")}, "Title")
	newModel, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	lm := newModel.(ListModel)

	if lm.List.Width() != 80 {
		t.Errorf("expected list width to be 80, got %d", lm.List.Width())
	}
	if lm.List.Height() != 24 {
		t.Errorf("expected list height to be 24, got %d", lm.List.Height())
	}
}

func TestView(t *testing.T) {
	items := []list.Item{
		NewItem("PostgreSQL", "PostgreSQL database"),
		NewItem("MySQL", "MySQL database"),
	}
	title := "Choose Database"
	m := CreateListModel(items, title)

	view := m.View()
	if !strings.Contains(view, title) {
		t.Errorf("expected view to contain title %q, got:\n%s", title, view)
	}
	if !strings.Contains(view, "PostgreSQL") {
		t.Errorf("expected view to contain PostgreSQL, got:\n%s", view)
	}

	m.Quitting = true
	if m.View() != "" {
		t.Errorf("expected empty view on quitting, got %q", m.View())
	}
}
