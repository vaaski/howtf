package main

import "strings"

func queryView(m *model) string {
	var s string

	s += "Query:"
	s += "\n\n"
	s += strings.Join(m.args, " ")
	s += "\n"

	return s
}
