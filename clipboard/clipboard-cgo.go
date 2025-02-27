//go:build cgo

package clipboard

import (
	"golang.design/x/clipboard"
)

var ClipboardAvailable = true

func WriteToClipboard(text string) {
	clipboard.Write(clipboard.FmtText, []byte(text))
}
