package metrics

import (
	"azul3d.org/engine/keyboard"
	"github.com/rikonor/keysig/keylogger"
)

// TimeToNext metric type keeps track of the time delta between different characters
// E.g. the time from keyUp of key A and keyDown of key B
type TimeToNext struct {
	inputChan chan keyboard.ButtonEvent
	active    bool
}

func NewTimeToNext() *TimeToNext {
	return &TimeToNext{
		inputChan: make(chan keyboard.ButtonEvent),
	}
}

func (m *TimeToNext) consumeStream() {
	for evt := range m.inputChan {
		m.processEvent(evt)
	}
}

func (m *TimeToNext) RegisterWith(k *keylogger.Keylogger) {
	k.Register("timeToNext", m.inputChan)

	if !m.active {
		go m.consumeStream()
		m.active = true
	}
}

// Implementation

func (m *TimeToNext) processEvent(evt keyboard.ButtonEvent) {

}
