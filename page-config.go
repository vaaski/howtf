package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zalando/go-keyring"
)

const (
	tokenInputFocus = iota
	modelInputFocus
)

type configModel struct {
	focusedInput int
	quitting     bool

	tokenInput textinput.Model
	modelInput textinput.Model

	openAIToken     string
	openAITextModel string
}

func configInitialModel() configModel {
	config := configModel{
		openAITextModel: "gpt-4o",
		openAIToken:     "",
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
	config.tokenInput.Prompt = "OpenAI API Token: "
	config.tokenInput.SetValue(config.openAIToken)

	config.modelInput = textinput.New()
	config.modelInput.Placeholder = "Enter OpenAI Text Model"
	config.modelInput.Prompt = "OpenAI Text Model: "
	config.modelInput.SetValue(config.openAITextModel)

	return config
}

func configView(m *model) string {
	var s string

	s += "CONFIG"
	s += "\n\n"

	if !m.config.quitting {
		s += m.config.tokenInput.View()
		s += "\n"
		s += m.config.modelInput.View()
	} else {
		s += "OpenAI API Token: "

		if len(m.config.openAIToken) > 0 {
			s += hideString(m.config.openAIToken)
		} else {
			s += "NOT SET"
		}

		s += "\n"

		s += "OpenAI Text Model: "
		s += m.config.openAITextModel
	}

	s += "\n"
	return s
}

func configController(m *model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			m.config.focusedInput = (m.config.focusedInput + 1) % 2

		case tea.KeyShiftTab:
			m.config.focusedInput = (m.config.focusedInput - 1 + 2) % 2

		case tea.KeyEnter:
			m.config.openAIToken = m.config.tokenInput.Value()
			m.config.tokenInput.Blur()

			m.config.openAITextModel = m.config.modelInput.Value()
			m.config.modelInput.Blur()

			m.config.quitting = true

			if len(m.config.openAIToken) > 0 {
				keyring.Set(SERVICE_NAME, TOKEN_NAME, m.config.openAIToken)
			}

			if len(m.config.openAITextModel) > 0 {
				keyring.Set(SERVICE_NAME, MODEL_NAME, m.config.openAITextModel)
			}

			return m, tea.Quit

		case tea.KeyEsc:
			m.config.quitting = true
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
