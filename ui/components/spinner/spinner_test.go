package spinner

import (
	"errors"
	"strings"
	"testing"
)

func TestCreateModel(t *testing.T) {
	steps := []Step{
		{
			Title:     "Step 1",
			DoneTitle: "Step 1 Done",
			Action:    func() error { return nil },
		},
	}
	m := CreateModel("Test Title", steps)

	if m.Title != "Test Title" {
		t.Errorf("expected Title 'Test Title', got %q", m.Title)
	}
	if len(m.Steps) != 1 {
		t.Errorf("expected 1 step, got %d", len(m.Steps))
	}
	if m.CurrentStep != 0 {
		t.Errorf("expected CurrentStep 0, got %d", m.CurrentStep)
	}
	if m.Done {
		t.Errorf("expected Done to be false")
	}
}

func TestStepExecutionSuccess(t *testing.T) {
	executed := false
	steps := []Step{
		{
			Title:     "Step 1",
			DoneTitle: "Step 1 Done",
			Action: func() error {
				executed = true
				return nil
			},
		},
	}

	m := CreateModel("Test Title", steps)
	newModel, cmd := m.Update(StepDoneMsg{Index: 0})
	sm := newModel.(SpinnerModel)

	if sm.CurrentStep != 1 {
		t.Errorf("expected CurrentStep 1, got %d", sm.CurrentStep)
	}
	if !sm.Done {
		t.Errorf("expected Done to be true after finishing all steps")
	}
	if cmd == nil {
		t.Errorf("expected tea.Quit command")
	}
	_ = executed
}

func TestStepExecutionError(t *testing.T) {
	testErr := errors.New("something went wrong")
	steps := []Step{
		{
			Title:  "Step 1",
			Action: func() error { return testErr },
		},
	}

	m := CreateModel("Test Title", steps)
	newModel, cmd := m.Update(StepErrMsg{Index: 0, Err: testErr})
	sm := newModel.(SpinnerModel)

	if !sm.Done {
		t.Errorf("expected Done to be true on error")
	}
	if sm.Err == nil || sm.Err.Error() != testErr.Error() {
		t.Errorf("expected error %v, got %v", testErr, sm.Err)
	}
	if cmd == nil {
		t.Errorf("expected tea.Quit command")
	}
}

func TestSpinnerView(t *testing.T) {
	steps := []Step{
		{
			Title:     "Creating folders...",
			DoneTitle: "Created folders",
			Action:    func() error { return nil },
		},
		{
			Title:     "Installing dependencies...",
			DoneTitle: "Installed dependencies",
			Action:    func() error { return nil },
		},
	}

	m := CreateModel("Generating Project", steps)
	view := m.View()

	if !strings.Contains(view, "Generating Project") {
		t.Errorf("expected view to contain title 'Generating Project', got:\n%s", view)
	}
	if !strings.Contains(view, "Creating folders...") {
		t.Errorf("expected view to contain step title, got:\n%s", view)
	}

	// Move to step 2
	m.CurrentStep = 1
	view = m.View()
	if !strings.Contains(view, "Created folders") {
		t.Errorf("expected view to contain done title 'Created folders', got:\n%s", view)
	}
}
