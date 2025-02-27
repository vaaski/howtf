//go:build !cgo

package clipboard

var ClipboardAvailable = false

func WriteToClipboard(text string) {
}
