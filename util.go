package main

import "strings"

func hideString(s string) string {
	return strings.Repeat("*", len(s))
}
