package keylogger

import (
	"fmt"

	"azul3d.org/engine/gfx"
	"azul3d.org/engine/gfx/window"
	"github.com/rikonor/keysig/keylogger/keyboard"
)

type Keylogger struct {
	consumers map[string]chan keyboard.ButtonEvent
}

func New() *Keylogger {
	return &Keylogger{
		consumers: make(map[string]chan keyboard.ButtonEvent),
	}
}

func (k *Keylogger) Register(cName string, cChan chan keyboard.ButtonEvent) {
	k.consumers[cName] = cChan
}

func (k *Keylogger) Deregister(cName string) {
	delete(k.consumers, cName)
}

// broadcastEvent takes an event and sends it to all registered consumers
func (k *Keylogger) broadcastEvent(e keyboard.ButtonEvent) {
	for _, c := range k.consumers {
		c <- e
	}
}

// startStreaming consumes button events from the global window
// and forwards them to the registered consumers
func (k *Keylogger) startStreaming(w window.Window) {
	events := make(chan window.Event, 256)
	w.Notify(events, window.KeyboardButtonEvents)

	go func() {
		for event := range events {
			// Try converting button event from Azul3D
			a3devt, ok := extractAzul3DButtonEvent(event)

			// Make sure no unrelated events get through
			if !ok {
				fmt.Println("Not a ButtonEvent event")
				continue
			}

			evt := fromAzul3DEvent(a3devt)

			// // Make sure no unrelated events get through
			// evt, ok := event.(keyboard.ButtonEvent)
			// if !ok {
			// 	fmt.Println("Not a ButtonEvent event")
			// 	continue
			// }

			k.broadcastEvent(evt)
		}
	}()
}

// blockForever renders a black screen and blocks forever
func (k *Keylogger) blockForever(d gfx.Device) {
	for {
		d.Render()
	}
}

// Start ...
func (k *Keylogger) Start() {
	window.Run(func(w window.Window, d gfx.Device) {
		k.startStreaming(w)
		k.blockForever(d)
	}, nil)
}
