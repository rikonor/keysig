package keylogger

import (
	azul3d "azul3d.org/engine/keyboard"
	"github.com/rikonor/keysig/keylogger/keyboard"
)

func extractAzul3DButtonEvent(be interface{}) (azul3d.ButtonEvent, bool) {
	evt, ok := be.(azul3d.ButtonEvent)
	return evt, ok
}

func isAzul3DButtonEvent(be interface{}) bool {
	_, ok := be.(azul3d.ButtonEvent)
	return ok
}

func fromAzul3DEvent(be azul3d.ButtonEvent) keyboard.ButtonEvent {
	return keyboard.ButtonEvent{
		T:     be.T,
		Key:   fromAzul3DKey(be.Key),
		State: fromAzul3DState(be.State),
		Raw:   be.Raw,
	}
}

func fromAzul3DKey(k azul3d.Key) keyboard.Key {
	return keyboard.Key(k)
}

func fromAzul3DState(s azul3d.State) keyboard.State {
	return keyboard.State(s)
}
