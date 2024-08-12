package main

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/wordwrap"
)

type queryModel struct {
	query    string
	response string

	loading         bool
	responseChannel chan queryResponse

	input textinput.Model
}

type queryResponse string
type responseFinished bool

func queryInitialModel() queryModel {
	query := queryModel{
		responseChannel: make(chan queryResponse),
	}

	query.input = textinput.New()
	query.input.Placeholder = "Enter query"
	query.input.Prompt = "Query: "

	return query
}

func queryView(m *model) string {
	var s string

	s += "QUERY"

	if m.query.loading {
		s += " (loading)"
	}

	s += "\n\n"
	if len(m.query.query) == 0 {
		s += m.query.input.View()
	} else {
		s += "Query: "
		s += m.query.query
	}

	s += "\n\n"
	s += wordwrap.String(m.query.response, m.termWidth)

	s += "\n"

	return s
}

func queryController(m *model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	if len(m.query.query) == 0 {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				log.Println("enter")
				m.query.loading = true
				m.query.input.Blur()
				m.query.query = m.query.input.Value()
				return m, tea.Batch(generateResponse(m), awaitResponse(m))
			}
		}
		m.query.input, cmd = m.query.input.Update(msg)
	}

	switch msg := msg.(type) {
	case initCmd:
		if len(m.args) == 0 {
			m.query.input.Focus()
		} else {
			m.query.query = strings.Join(m.args, " ")
			m.query.loading = true
			return m, tea.Batch(generateResponse(m), awaitResponse(m))
		}

	case queryResponse:
		if m.query.loading {
			log.Println("queryResponse", msg)
			m.query.response = string(msg)
			return m, awaitResponse(m)
		}

	case responseFinished:
		log.Println("responseFinished", msg)
		m.query.loading = false
		close(m.query.responseChannel)
		return m, nil
	}

	return m, cmd
}

func generateResponse(m *model) tea.Cmd {
	return func() tea.Msg {
		generateGPT(m.config.openAIToken, m.query.query, m.query.responseChannel)
		return responseFinished(true)
	}
}

func awaitResponse(m *model) tea.Cmd {
	return func() tea.Msg {
		return queryResponse(<-m.query.responseChannel)
	}
}
