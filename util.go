package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
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
