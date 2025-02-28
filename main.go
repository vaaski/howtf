package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vaaski/howtf/executor"
)

var originalQuery string
var responseMarkdown string
var shouldExecute = false
var explainMode = false

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

	if len(responseMarkdown) > 0 {
		out, err := markdownCodeRenderer.Render(responseMarkdown)
		if err != nil {
			log.Println("error rendering markdown", err)
		}

		if explainMode {
			fmt.Println(borderStyle.UnsetWidth().Padding(1, 2, 0, 0).Render(greyedOutStyle.Render(out)))
		} else {
			fmt.Println(borderStyle.UnsetWidth().Padding(0).Render(greyedOutStyle.Render(out)))
		}

		if shouldExecute {
			commandToExecute := extractMarkdownMaybe(responseMarkdown)
			executor.Execute(commandToExecute)
		}
	}
}
