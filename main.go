package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vulan1999/go-gin-template-cli/ui/components/choice"
	inputtext "github.com/vulan1999/go-gin-template-cli/ui/components/input_text"
)

func main() {
	p := tea.NewProgram(inputtext.CreateModel("Enter your Gin Project Name:", "Your Gin Project Name", true))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
	c := tea.NewProgram(choice.CreateChoicesModel([]string{
		"Barebone (Local only -> e.q., myapp)",
		"Github (Remote Github repo -> e.g., github.com/username/myapp)",
		"Gitlab (Remote Gitlab repo -> e.g., gitlab.com/username/myapp)",
	}, "Choose your module convention"))

	if _, err := c.Run(); err != nil {
		log.Fatal(err)
	}
}
