package keylogger

import (
	"fmt"

	"azul3d.org/engine/gfx"
	"azul3d.org/engine/gfx/window"

	azul3d "azul3d.org/engine/keyboard"
	"github.com/rikonor/keysig/keylogger/keyboard"
)

// Azul3D keylogger

type Keylogger struct {
	outputChannel *chan keyboard.ButtonEvent
}

func NewKeylogger(outputChannel *chan keyboard.ButtonEvent) *Keylogger {
	return &Keylogger{outputChannel: outputChannel}
}

// startStreaming consumes button events from the global window
// and forwards them to the registered consumers
func (k *Keylogger) startStreaming(w window.Window) {
	events := make(chan window.Event, 256)
	w.Notify(events, window.KeyboardButtonEvents)

	go func() {
		for event := range events {
			// Try converting button event from Azul3D
			// and make sure no unrelated events get through
			a3devt, ok := extractAzul3DButtonEvent(event)
			if !ok {
				fmt.Println("Not a ButtonEvent event")
				continue
			}

			evt := fromAzul3DEvent(a3devt)
			*k.outputChannel <- evt
		}
	}()
}

// blockForever renders a black screen and blocks forever
func (k *Keylogger) blockForever(d gfx.Device) {
	for {
		d.Render()
	}
}

func (k *Keylogger) Start() {
	window.Run(func(w window.Window, d gfx.Device) {
		k.startStreaming(w)
		k.blockForever(d)
	}, nil)
}

// Utils

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
	}
}

func fromAzul3DKey(k azul3d.Key) keyboard.Key {
	return keyboard.Key(k)
}

func fromAzul3DState(s azul3d.State) keyboard.State {
	return keyboard.State(s)
}
