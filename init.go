package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// love go enums
const (
	ConfigPage = iota
	GenPage
)

type model struct {
	args []string
	page int
}

func initialModel() model {
	return model{
		args: os.Args[1:],
		page: GenPage,
	}
}

func (m model) Init() tea.Cmd {

	if m.page == GenPage {
		return tea.Quit
	}

	return nil
}
