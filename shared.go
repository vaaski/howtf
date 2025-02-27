package main

import "github.com/charmbracelet/lipgloss"

// style definitions
var (
	borderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#ad0c88")).
			Padding(1, 2)

	borderWidth = lipgloss.Width(borderStyle.Render(""))

	headerString = lipgloss.NewStyle().
			Padding(0, 3).
			SetString("howtf").
			String()

	headerStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("#ffffff"))

	chevronStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#0ff"))
	greyedOutStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
)
