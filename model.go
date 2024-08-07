package main

import (
	"flag"

	tea "github.com/charmbracelet/bubbletea"
)

// self incrementing page constants
const (
	QueryPage = iota
	ConfigPage
)

type setPage int

const (
	SERVICE_NAME = "howtf-cli"
	TOKEN_NAME   = "OpenAI-API-Token"
	MODEL_NAME   = "OpenAI-Model"
)

type flags struct {
	config *bool
}

type model struct {
	args []string
	page int

	config configModel
	query  queryModel

	flags flags
}

func initialModel() model {
	m := model{
		page:   QueryPage,
		config: configInitialModel(),
		query:  queryInitialModel(),

		flags: flags{
			config: flag.Bool("config", false, "Open config page"),
		},
	}

	flag.Parse()
	m.args = flag.Args()

	return m
}

type initCmd struct{}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			if *m.flags.config {
				return setPage(ConfigPage)
			}

			return setPage(QueryPage)
		},
		func() tea.Msg {
			return initCmd{}
		},
	)
}
