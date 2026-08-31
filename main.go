package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vulan1999/go-gin-template-cli/ui"
)

func main() {
	p := tea.NewProgram(ui.NewWizard())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
