package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zalando/go-keyring"
)

const (
	tokenInputFocus = iota
	modelInputFocus
)

type configKeyMap struct {
	Help     key.Binding
	Quit     key.Binding
	Save     key.Binding
	Next     key.Binding
	Previous key.Binding
	Delete   key.Binding
}

func (k configKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Save, k.Help}
}

func (k configKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Save, k.Next, k.Previous, k.Delete},
		{k.Help, k.Quit},
	}
}

var configKeys = configKeyMap{
	Help: key.NewBinding(
		key.WithKeys("ctrl+h"),
		key.WithHelp("ctrl+h", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc", "cancel"),
	),
	Save: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "save"),
	),
	Next: key.NewBinding(
		key.WithKeys("tab", "down"),
		key.WithHelp("↓/tab", "next input"),
	),
	Previous: key.NewBinding(
		key.WithKeys("shift+tab", "up"),
		key.WithHelp("↑/shift+tab", "previous input"),
	),
	Delete: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "delete data"),
	),
}

type configModel struct {
	focusedInput int

	tokenInput textinput.Model
	modelInput textinput.Model

	openAIToken     string
	openAITextModel string

	keys configKeyMap
	help help.Model
}

func configInitialModel() configModel {
	config := configModel{
		openAITextModel: "gpt-4o",
		openAIToken:     "",
		keys:            configKeys,
		help:            help.New(),
	}

	token, err := keyring.Get(SERVICE_NAME, TOKEN_NAME)
	if err == nil {
		config.openAIToken = token
	}

	model, err := keyring.Get(SERVICE_NAME, MODEL_NAME)
	if err == nil {
		config.openAITextModel = model
	}

	config.tokenInput = textinput.New()
	config.tokenInput.Placeholder = "Enter OpenAI API Token"
	config.tokenInput.Prompt = "OpenAI API Token   "
	config.tokenInput.PromptStyle = greyedOutStyle
	config.tokenInput.SetValue(config.openAIToken)

	config.modelInput = textinput.New()
	config.modelInput.Placeholder = "Enter OpenAI Text Model"
	config.modelInput.Prompt = "OpenAI Text Model  "
	config.modelInput.PromptStyle = greyedOutStyle
	config.modelInput.SetValue(config.openAITextModel)

	return config
}

func configView(m *model) string {
	var s string

	headerPositioning := lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Margin(1, 0)

	configHeaderString := lipgloss.NewStyle().
		Padding(0, 3).
		SetString("howtf config").
		String()

	s += headerPositioning.Render(
		gradientBackgroundText(headerStyle, configHeaderString, lipgloss.Color("#ad0c88"), lipgloss.Color("#d65ab9")),
	)
	s += "\n"

	s += "\n"

	var inputs string

	inputs += m.config.tokenInput.View()
	inputs += "\n"
	inputs += m.config.modelInput.View()

	s += borderStyle.Render(inputs)

	s += "\n\n"
	s += m.config.help.View(m.config.keys)

	return s
}

func configController(m *model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, configKeys.Next):
			m.config.focusedInput = (m.config.focusedInput + 1) % 2

		case key.Matches(msg, configKeys.Previous):
			m.config.focusedInput = (m.config.focusedInput - 1 + 2) % 2

		case key.Matches(msg, configKeys.Save):
			m.config.openAIToken = m.config.tokenInput.Value()
			m.config.tokenInput.Blur()

			m.config.openAITextModel = m.config.modelInput.Value()
			m.config.modelInput.Blur()

			if len(m.config.openAIToken) > 0 {
				keyring.Set(SERVICE_NAME, TOKEN_NAME, m.config.openAIToken)
			}

			if len(m.config.openAITextModel) > 0 {
				keyring.Set(SERVICE_NAME, MODEL_NAME, m.config.openAITextModel)
			}

			return m, tea.Quit

		case key.Matches(msg, configKeys.Quit):
			return m, tea.Quit

		case key.Matches(msg, configKeys.Help):
			m.config.help.ShowAll = !m.config.help.ShowAll
			return m, nil

		case key.Matches(msg, configKeys.Delete):
			keyring.Delete(SERVICE_NAME, TOKEN_NAME)
			keyring.Delete(SERVICE_NAME, MODEL_NAME)
			return m, tea.Quit
		}

	}

	m.config.tokenInput.Blur()
	m.config.modelInput.Blur()

	switch m.config.focusedInput {
	case tokenInputFocus:
		m.config.tokenInput.Focus()
		m.config.tokenInput, cmd = m.config.tokenInput.Update(msg)

	case modelInputFocus:
		m.config.modelInput.Focus()
		m.config.modelInput, cmd = m.config.modelInput.Update(msg)
	}

	return m, cmd
}
