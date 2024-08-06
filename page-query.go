package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type queryModel struct {
	initialized bool
	responseRaw string

	input textinput.Model
}

func queryInitialModel() queryModel {
	query := queryModel{}

	query.input = textinput.New()
	query.input.Placeholder = "Enter query"
	query.input.Prompt = ": "

	return query
}

func queryView(m *model) string {
	var s string

	s += "QUERY"
	s += "\n\n"
	s += strings.Join(m.args, " ")
	s += m.query.input.View()
	s += "\n"

	return s
}

func queryController(m *model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if !m.query.initialized {
		if len(m.args) == 0 {
			m.args = append(m.args, "should ask for input")
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				return m, nil

			case tea.KeyEsc:
				return m, tea.Quit
			}
		}

		m.query.input.Focus()
		m.query.input, cmd = m.query.input.Update(msg)
	}

	return m, cmd
}
