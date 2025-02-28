package main

import (
	"image/color"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/glamour/ansi"
	"github.com/charmbracelet/glamour/styles"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
)

func getUserShell() string {
	var path string

	switch runtime.GOOS {
	case "windows":
		if os.Getenv("COMSPEC") != "" {
			path = os.Getenv("COMSPEC")
		} else {
			path = "cmd.exe"
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

var MARKDOWN_REGEX = regexp.MustCompile(`(?s)\x60\x60\x60(?:\w*)?(?:\n)(.+)(?:\n)\x60\x60\x60|\x60(.+)\x60`)

func extractMarkdownMaybe(s string) string {
	matches := MARKDOWN_REGEX.FindStringSubmatch(s)

	if len(matches) <= 1 {
		return s
	} else if len(matches[1]) > 0 {
		return matches[1]
	} else if len(matches[2]) > 0 {
		return matches[2]
	}

	return s
}

func uintPtr(u uint) *uint { return &u }

var markdownCodeRenderer, _ = glamour.NewTermRenderer(
	glamour.WithStyles(func() ansi.StyleConfig {
		styles := styles.DarkStyleConfig
		styles.Document.BlockPrefix = ""
		styles.Document.BlockSuffix = ""
		styles.CodeBlock.Margin = uintPtr(0)
		return styles
	}()),
)
