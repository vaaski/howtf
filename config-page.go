package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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
	}

	config.tokenInput = textinput.New()
	config.tokenInput.Placeholder = "Enter OpenAI API Token"
	config.tokenInput.Prompt = ": "
	config.modelInput.SetValue(config.openAIToken)

	config.modelInput = textinput.New()
	config.modelInput.Placeholder = "Enter OpenAI Text Model"
	config.modelInput.Prompt = ": "
	config.modelInput.SetValue(config.openAITextModel)

	return config
}

func configView(m *model) string {
	var s string

	s += "CONFIG"
	s += "\n\n"

	if m.config.quitting {
		s += "OpenAI API Token: "

		if len(m.config.openAIToken) > 0 {
			s += hideString(m.config.openAIToken)
		} else {
			s += "NOT SET"
		}

		s += "\n"

		s += "OpenAI Text Model: "
		s += m.config.openAITextModel
	} else {
		s += "Token"
		s += m.config.tokenInput.View()

		s += "\n"

		s += "Model"
		s += m.config.modelInput.View()
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
