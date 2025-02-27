package main

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
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

// style definitions
var (
	borderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#ad0c88")).
			Padding(1, 2)

	borderWidth = lipgloss.Width(borderStyle.Render(""))

	headerString = lipgloss.NewStyle().
			Padding(0, 2).
			SetString("howtf").
			String()

	headerStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("#ffffff"))

	chevronStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#0ff"))
	greyedOutStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
)

func queryInitialModel() queryModel {
	query := queryModel{
		responseChannel: make(chan queryResponse),
	}

	query.queryInput = textinput.New()
	query.queryInput.Placeholder = "Describe the problem you want to solve"
	query.queryInput.Prompt = chevronStyle.Render("> ")

	return query
}

func queryView(m *model) string {
	var s string

	borderStyle = borderStyle.Width(m.termWidth - borderWidth)
	headerPositioning := lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Margin(1, 0)

	s += headerPositioning.Render(gradientBackgroundText(headerStyle, headerString, lipgloss.Color("#ad0c88"), lipgloss.Color("#d65ab9")))

	s += "\n"
	if len(m.query.finalQuery) == 0 {
		s += borderStyle.Render(m.query.queryInput.View())
	} else {
		queryStyle := borderStyle.UnsetBorderStyle().Padding(0, 2)
		s += queryStyle.Render(chevronStyle.Render("> ") + greyedOutStyle.Render(m.query.finalQuery))
	}

	if len(m.query.response) > 0 {
		s += "\n"
		s += borderStyle.Render(wordwrap.String(m.query.response, m.termWidth))
	}

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

	// the query is empty, handle text input events
	if len(m.query.finalQuery) == 0 {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				queryInput := m.query.queryInput.Value()

				// ignore enter key on empty input
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
				// only show the original query if it was entered interactively
				if len(m.args) == 0 {
					originalQuery = m.query.finalQuery
				}

				commandToExecute = m.query.response
				return m, tea.Quit
			case "n":
				return m, tea.Quit
			}
		}
	}

	switch msg := msg.(type) {
	case initCmd:
		// if query is given as arguments, skip the text input
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
