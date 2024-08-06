package main

func (m model) View() string {
	var s string

	if m.page == QueryPage {
		s += queryView(&m)
	} else if m.page == ConfigPage {
		s += configView(&m)
	}

	return s
}
