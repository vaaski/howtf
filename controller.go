package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height

	case setPage:
		m.page = int(msg)

	case tea.KeyMsg:
		switch msg.String() {
		}
	}

	switch m.page {
	case ConfigPage:
		return configController(&m, msg)
	case QueryPage:
		return queryController(&m, msg)
	}

	return m, tea.Batch(cmd)
}
