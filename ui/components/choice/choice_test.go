package choice

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestCreateChoicesModel(t *testing.T) {
	choices := []string{"Option 1", "Option 2", "Option 3"}
	title := "Select an option"
	m := CreateChoicesModel(choices, title)

	if m.Title != title {
		t.Errorf("expected Title %q, got %q", title, m.Title)
	}
	if len(m.Choices) != len(choices) {
		t.Errorf("expected %d choices, got %d", len(choices), len(m.Choices))
	}
	if m.Cursor != 0 {
		t.Errorf("expected Cursor to start at 0, got %d", m.Cursor)
	}
	if m.Quitting {
		t.Errorf("expected Quitting to be false")
	}
}

func TestInit(t *testing.T) {
	m := CreateChoicesModel([]string{"A"}, "Title")
	cmd := m.Init()
	if cmd != nil {
		t.Errorf("expected Init to return nil, got %v", cmd)
	}
}

func TestUpdate_CtrlCAndEsc(t *testing.T) {
	tests := []struct {
		name    string
		keyMsg  tea.KeyMsg
	}{
		{"Ctrl+C by KeyMsg", tea.KeyMsg{Type: tea.KeyCtrlC}},
		{"Esc by KeyMsg", tea.KeyMsg{Type: tea.KeyEsc}},
		{"ctrl+c by string", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}, Alt: false}},
	}

	for _, tt := range tests[:2] {
		t.Run(tt.name, func(t *testing.T) {
			m := CreateChoicesModel([]string{"A", "B"}, "Title")
			newModel, cmd := m.Update(tt.keyMsg)
			cm, ok := newModel.(ChoicesModel)
			if !ok {
				t.Fatalf("expected ChoicesModel, got %T", newModel)
			}
			if !cm.Quitting {
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

func TestUpdate_NavigationAndSelection(t *testing.T) {
	choices := []string{"First", "Second", "Third"}
	m := CreateChoicesModel(choices, "Title")

	// Navigate down with "down"
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	cm := newModel.(ChoicesModel)
	if cm.Cursor != 1 {
		t.Errorf("expected cursor 1 after down key, got %d", cm.Cursor)
	}

	// Navigate down with "j"
	newModel, _ = cm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	cm = newModel.(ChoicesModel)
	if cm.Cursor != 2 {
		t.Errorf("expected cursor 2 after 'j', got %d", cm.Cursor)
	}

	// Cannot navigate past the last choice
	newModel, _ = cm.Update(tea.KeyMsg{Type: tea.KeyDown})
	cm = newModel.(ChoicesModel)
	if cm.Cursor != 2 {
		t.Errorf("expected cursor to remain 2 at bottom, got %d", cm.Cursor)
	}

	// Navigate up with "up"
	newModel, _ = cm.Update(tea.KeyMsg{Type: tea.KeyUp})
	cm = newModel.(ChoicesModel)
	if cm.Cursor != 1 {
		t.Errorf("expected cursor 1 after up key, got %d", cm.Cursor)
	}

	// Navigate up with "k"
	newModel, _ = cm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	cm = newModel.(ChoicesModel)
	if cm.Cursor != 0 {
		t.Errorf("expected cursor 0 after 'k', got %d", cm.Cursor)
	}

	// Cannot navigate above the first choice
	newModel, _ = cm.Update(tea.KeyMsg{Type: tea.KeyUp})
	cm = newModel.(ChoicesModel)
	if cm.Cursor != 0 {
		t.Errorf("expected cursor to remain 0 at top, got %d", cm.Cursor)
	}

	// Select with "enter"
	newModel, _ = cm.Update(tea.KeyMsg{Type: tea.KeyEnter})
	cm = newModel.(ChoicesModel)
	if !cm.Done {
		t.Errorf("expected Done to be true on enter")
	}
	if cm.Choice != 0 {
		t.Errorf("expected Choice to be 0, got %d", cm.Choice)
	}
}

func TestView(t *testing.T) {
	choices := []string{"Alpha", "Beta"}
	title := "Choose Greek Letter"
	m := CreateChoicesModel(choices, title)

	view := m.View()

	if !strings.Contains(view, title) {
		t.Errorf("expected view to contain title %q", title)
	}
	if !strings.Contains(view, "(*) Alpha") {
		t.Errorf("expected view to indicate selected cursor at Alpha, got:\n%s", view)
	}
	if !strings.Contains(view, "( ) Beta") {
		t.Errorf("expected view to indicate unselected Beta, got:\n%s", view)
	}
}
