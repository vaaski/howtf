package main

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/wordwrap"
	"github.com/vaaski/howtf/executor"
)

type queryModel struct {
	finalQuery string
	response   string

	loading         bool
	finished        bool
	responseChannel chan queryResponse

	queryInput textinput.Model
}

type queryResponse string
type responseFinished bool

func queryInitialModel() queryModel {
	query := queryModel{
		responseChannel: make(chan queryResponse),
	}

	query.queryInput = textinput.New()
	query.queryInput.Placeholder = "Describe the problem you want to solve"
	query.queryInput.Prompt = "> "

	return query
}

func queryView(m *model) string {
	var s string

	s += "QUERY PAGE"

	if m.query.loading {
		s += " (loading)"
	}

	s += "\n\n"
	if len(m.query.finalQuery) == 0 {
		s += m.query.queryInput.View()
	} else {
		s += "Query: "
		s += m.query.finalQuery
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

	if len(m.query.finalQuery) == 0 {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				queryInput := m.query.queryInput.Value()
				if len(queryInput) == 0 {
					return m, nil
				}

				log.Println("enter")
				m.query.loading = true
				m.query.queryInput.Blur()
				m.query.finalQuery = queryInput
				return m, tea.Batch(generateResponse(m), awaitResponse(m))
			}
		}
		m.query.queryInput, cmd = m.query.queryInput.Update(msg)
	}

	if m.query.loading == false && m.query.finished == true {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "y", "enter":
				executor.Execute(m.query.response)
				return m, tea.Quit
			case "n":
				return m, tea.Quit
			}
		}
	}

	switch msg := msg.(type) {
	case initCmd:
		if len(m.args) == 0 {
			m.query.queryInput.Focus()
			return m, textinput.Blink
		} else {
			m.query.finalQuery = strings.Join(m.args, " ")
			m.query.loading = true
			return m, tea.Batch(generateResponse(m), awaitResponse(m))
		}

	case responseFinished:
		log.Println("responseFinished", msg)
		m.query.loading = false
		m.query.finished = true
		close(m.query.responseChannel)
		return m, nil

	case queryResponse:
		if m.query.loading {
			log.Println("queryResponse", msg)
			m.query.response = string(msg)
			return m, awaitResponse(m)
		}
	}

	return m, cmd
}

func generateResponse(m *model) tea.Cmd {
	return func() tea.Msg {
		generateGPT(m.config.openAIToken, m.query.finalQuery, m.query.responseChannel)
		return responseFinished(true)
	}
}

func awaitResponse(m *model) tea.Cmd {
	return func() tea.Msg {
		return queryResponse(<-m.query.responseChannel)
	}
}
