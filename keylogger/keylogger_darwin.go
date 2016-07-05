package keylogger

import (
	"github.com/rikonor/keysig/keylogger/keyboard"
	"github.com/rikonor/keysig/keylogger/osxkeylogger"
)

func NewKeylogger(outputChannel *chan keyboard.ButtonEvent) *osxkeylogger.Keylogger {
	return osxkeylogger.NewKeylogger(outputChannel)
}
