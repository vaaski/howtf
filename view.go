package main

import (
	"strings"
)

func (m model) View() string {
	s := ""

	if m.page == GenPage {
		s += "Query:\n\n"
		s += strings.Join(m.args, " ")
		s += "\n"
	} else if m.page == ConfigPage {
		s += "Config:\n\n"
		s += "Coming soon!"
	}

	return s
}
