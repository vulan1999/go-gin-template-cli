package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vulan1999/go-gin-template-cli/ui"
)

func main() {
	wizard := ui.NewWizard()
	p := tea.NewProgram(wizard)
	model, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

	finalWizard, ok := model.(ui.WizardModel)
	if !ok || finalWizard.Cancelled {
		log.Println("Operation Cancelled")
		return
	}

	if finalWizard.Err != nil {
		log.Fatalf("Project generation failed: %v", finalWizard.Err)
	}
}
