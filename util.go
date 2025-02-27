package main

import (
	"image/color"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
)

func hideString(s string) string {
	return strings.Repeat("*", len(s))
}

func getUserShell() string {
	var path string

	switch runtime.GOOS {
	case "windows":
		if os.Getenv("COMSPEC") != "" {
			path = os.Getenv("COMSPEC")
		} else {
			path = "/cmd.exe"
		}
	case "darwin":
		if os.Getenv("SHELL") != "" {
			path = os.Getenv("SHELL")
		} else {
			path = "/bin/bash"
		}
	default:
		if os.Getenv("SHELL") != "" {
			path = os.Getenv("SHELL")
		} else {
			path = "/bin/sh"
		}
	}

	if path == "" {
		return "default shell"
	} else {
		return filepath.Base(path)
	}
}

func gradientBackgroundText(base lipgloss.Style, s string, c1 color.Color, c2 color.Color) string {
	colors := gamut.Blends(c1, c2, len(s))

	var str string
	for i, ss := range s {
		color, _ := colorful.MakeColor(colors[i%len(colors)])
		str = str + base.Background(lipgloss.Color(color.Hex())).Render(string(ss))
	}
	return str
}
