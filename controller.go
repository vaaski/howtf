package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case setPage:
		m.page = int(msg)

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	switch m.page {
	case ConfigPage:
		return configController(&m, msg)
	}

	return m, tea.Batch(cmd)
}
