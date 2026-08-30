package inputtext

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestTextInputValidation(t *testing.T) {
	// For required input case
	m := CreateModel("Enter Your Project Name", "Your Project Name", true)
	// Check input if user press enter right away
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	res := updated.(TextInputModel)

	if res.Error == nil || res.Quitting {
		t.Error("Expected Error, got nil")
	}

	// Check when user press Esc
	updatedEsc, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	resEsc := updatedEsc.(TextInputModel)
	if !resEsc.Quitting {
		t.Error("Expected quitting cli when pressing Esc")
	}

	// Check when user input valid string
	res.TextInput.SetValue("Test")
	updated, _ = res.Update(tea.KeyMsg{Type: tea.KeyEnter})
	res = updated.(TextInputModel)

	if !res.Quitting || res.Value() != "Test" {
		t.Error("Expected Quitting is true or project name is Test")
	}

	// For non-required input case
	mNonRequired := CreateModel("Enter sample", "Enter your sample", false)

	updatedNonRequired, _ := mNonRequired.Update(tea.KeyMsg{Type: tea.KeyEnter})
	resNonRequired := updatedNonRequired.(TextInputModel)
	if resNonRequired.Error != nil || !resNonRequired.Quitting {
		t.Errorf("Expected got no Error, got %v", resNonRequired.Error)
	}
}
