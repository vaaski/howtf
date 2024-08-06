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
		m.query.initialized = true

		if len(m.args) == 0 {
			m.args = append(m.args, "should ask for input")
		}
	}

	return m, nil
}
