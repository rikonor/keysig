package metrics

import (
	"fmt"

	"azul3d.org/engine/keyboard"
	"github.com/rikonor/keysig/keylogger"
)

// CharDiff metric type keeps track of the time delta between different characters
// E.g. the time from keyUp of key A and keyDown of key B

type CharDiff struct {
	inputChan chan keyboard.ButtonEvent
	active    bool
}

func NewCharDiff() *CharDiff {
	return &CharDiff{
		inputChan: make(chan keyboard.ButtonEvent),
	}
}

func (m *CharDiff) processEvent(evt keyboard.ButtonEvent) {
	fmt.Println("Processing event", evt)
}

func (m *CharDiff) consumeStream() {
	for evt := range m.inputChan {
		m.processEvent(evt)
	}
}

func (m *CharDiff) RegisterWith(k *keylogger.Keylogger) {
	k.Register("charDiff", m.inputChan)

	if !m.active {
		go m.consumeStream()
		m.active = true
	}
}
