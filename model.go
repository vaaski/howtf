package main

import (
	"flag"

	tea "github.com/charmbracelet/bubbletea"
)

// self incrementing page constants
const (
	GenPage = iota
	ConfigPage
)

type setPage int

const (
	SERVICE_NAME = "howtf"
	TOKEN_NAME   = "token"
)

type flags struct {
	config *bool
}

type model struct {
	args []string
	page int

	config configModel

	flags flags
}

func initialModel() model {
	m := model{
		page:   GenPage,
		config: configInitialModel(),

		flags: flags{
			config: flag.Bool("config", false, "Open config page"),
		},
	}

	flag.Parse()
	m.args = flag.Args()

	return m
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg {

		// _, err := keyring.Get(SERVICE_NAME, TOKEN_NAME)

		if *m.flags.config {
			return setPage(ConfigPage)
		}

		return setPage(GenPage)
	}
}
