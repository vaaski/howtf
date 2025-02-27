package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vaaski/howtf/executor"
)

var originalQuery string
var commandToExecute string

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	model := initialModel()
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	if len(originalQuery) > 0 {
		fmt.Println(chevronStyle.Render("> ") + originalQuery)
	}

	if len(commandToExecute) > 0 {
		fmt.Println(borderStyle.UnsetWidth().Render(greyedOutStyle.Render(commandToExecute)))

		executor.Execute(commandToExecute)
	}
}
