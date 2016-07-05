package keylogger

import (
	"github.com/rikonor/keysig/keylogger/keyboard"
	"github.com/rikonor/keysig/keylogger/linuxkeylogger"
)

func NewKeylogger(outputChannel *chan keyboard.ButtonEvent) *linuxkeylogger.Keylogger {
	return linuxkeylogger.NewKeylogger(outputChannel)
}
