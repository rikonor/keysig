package keylogger

import (
	"github.com/rikonor/keysig/keylogger/azul3dkeylogger"
	"github.com/rikonor/keysig/keylogger/keyboard"
)

func NewKeylogger(outputChannel *chan keyboard.ButtonEvent) *azul3dkeylogger.Keylogger {
	return azul3dkeylogger.NewKeylogger(outputChannel)
}
