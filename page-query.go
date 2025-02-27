package main

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vaaski/howtf/clipboard"
)

type queryKeyMap struct {
	Help    key.Binding
	Quit    key.Binding
	Execute key.Binding
	Copy    key.Binding
	Edit    key.Binding
}

func (k queryKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Execute, k.Help}
}

func (k queryKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Execute, k.Copy, k.Edit},
		{k.Help, k.Quit},
	}
}

var queryKeys = queryKeyMap{
	Help: key.NewBinding(
		key.WithKeys("ctrl+h"),
		key.WithHelp("ctrl+h", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc", "quit"),
	),
	Execute: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "query"),
	),
	Copy: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "copy"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
}

type queryModel struct {
	finalQuery string
	response   string

	loading         bool
	finished        bool
	responseChannel chan queryResponse

	queryInput textinput.Model
	spinner    spinner.Model
	help       help.Model
	keys       queryKeyMap
}

type queryResponse string
type responseFinished bool

func queryInitialModel() queryModel {
	query := queryModel{
		responseChannel: make(chan queryResponse),
		keys:            queryKeys,
		help:            help.New(),
	}

	query.queryInput = textinput.New()
	query.queryInput.Placeholder = "Describe the problem you want to solve"
	query.queryInput.Prompt = chevronStyle.Render("> ")

	query.spinner.Spinner = spinner.Dot
	query.spinner.Style = chevronStyle

	query.keys.Copy.SetEnabled(false)
	query.keys.Edit.SetEnabled(false)

	return query
}

func queryView(m *model) string {
	var s string

	borderStyle = borderStyle.Width(m.termWidth - borderWidth)
	headerPositioning := lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Margin(1, 0)

	s += headerPositioning.Render(
		gradientBackgroundText(headerStyle, headerString, lipgloss.Color("#ad0c88"), lipgloss.Color("#d65ab9")),
	)

	s += "\n"
	if len(m.query.finalQuery) == 0 {
		s += borderStyle.Render(m.query.queryInput.View())
	} else {
		queryStyle := borderStyle.UnsetBorderStyle().Padding(0, 2)
		prefix := chevronStyle.Render(">")
		if m.query.loading {
			prefix = m.query.spinner.View()
		}

		s += queryStyle.Render(prefix + " " + greyedOutStyle.Render(m.query.finalQuery))
	}

	if len(m.query.response) > 0 {
		s += "\n"
		out, err := markdownRenderer.Render(m.query.response)
		if err != nil {
			log.Println("error rendering markdown", err)
			out = m.query.response
		}

		s += borderStyle.Padding(0).Render(out)
	}

	s += "\n\n"
	s += m.query.help.View(m.query.keys)

	return s
}

func queryController(m *model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.query.loading {
		m.query.spinner, cmd = m.query.spinner.Update(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, queryKeys.Quit):
			return m, tea.Quit
		case key.Matches(msg, queryKeys.Help):
			m.query.help.ShowAll = !m.query.help.ShowAll
			return m, nil
		}
	}

	// the query is empty, handle text input events
	if len(m.query.finalQuery) == 0 {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, queryKeys.Execute):
				queryInput := m.query.queryInput.Value()

				// ignore enter key on empty input
				if len(queryInput) == 0 {
					return m, nil
				}

				log.Println("enter")
				m.query.loading = true
				m.query.queryInput.Blur()
				m.query.finalQuery = queryInput
				return m, tea.Batch(generateResponse(m), awaitResponseChunk(m), m.query.spinner.Tick)
			}
		}
		m.query.queryInput, cmd = m.query.queryInput.Update(msg)
	}

	if m.query.loading == false && m.query.finished == true {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, queryKeys.Execute):
				shouldExecute = true
				return m, tea.Quit
			case key.Matches(msg, queryKeys.Copy):
				clipboard.WriteToClipboard(extractMarkdownMaybe(m.query.response))
				return m, tea.Quit
			case key.Matches(msg, queryKeys.Quit):
				return m, tea.Quit
			}
		}
	}

	switch msg := msg.(type) {
	case initCmd:
		m.query.help.Width = m.termWidth

		// if query is given as arguments, skip the text input
		if len(m.args) == 0 {
			m.query.queryInput.Focus()
			return m, textinput.Blink
		} else {
			m.query.finalQuery = strings.Join(m.args, " ")
			m.query.loading = true
			return m, tea.Batch(generateResponse(m), awaitResponseChunk(m), m.query.spinner.Tick)
		}

	case responseFinished:
		log.Println("responseFinished", msg)
		m.query.loading = false
		m.query.finished = true
		close(m.query.responseChannel)

		// only show the original query if it was entered interactively
		if len(m.args) == 0 {
			originalQuery = m.query.finalQuery
		}

		responseMarkdown = m.query.response

		m.query.keys.Execute.SetHelp(m.query.keys.Execute.Help().Key, "execute")
		if clipboard.ClipboardAvailable {
			m.query.keys.Copy.SetEnabled(true)
		}

		// todo: edit mode
		// m.query.keys.Edit.SetEnabled(true)

		return m, nil

	case queryResponse:
		if m.query.loading {
			log.Println("queryResponse", msg)
			m.query.response = string(msg)
			return m, awaitResponseChunk(m)
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

func awaitResponseChunk(m *model) tea.Cmd {
	return func() tea.Msg {
		return queryResponse(<-m.query.responseChannel)
	}
}
