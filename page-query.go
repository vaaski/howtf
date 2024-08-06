package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type queryModel struct {
	initialized bool
	responseRaw string
}

func queryInitialModel() queryModel {
	return queryModel{}
}

func queryView(m *model) string {
	var s string

	s += "Query:"
	s += "\n\n"
	s += strings.Join(m.args, " ")
	s += "\n"

	return s
}

func queryController(m *model, msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.query.initialized {
		if len(m.args) == 0 {
			m.args = append(m.args, "should ask for input")
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				m.config.openAIToken = m.config.tokenInput.Value()
				m.config.tokenInput.Blur()

				return m, nil

			case tea.KeyEsc:
				return m, tea.Quit
			}
		}
	}

	return m, nil
}
