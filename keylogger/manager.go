package keylogger

import (
	"github.com/rikonor/keysig/keylogger/keyboard"
)

type Manager struct {
	consumers map[string]chan keyboard.ButtonEvent
}

func NewManager() *Manager {
	return &Manager{
		consumers: make(map[string]chan keyboard.ButtonEvent),
	}
}

func (m *Manager) Register(cName string, cChan chan keyboard.ButtonEvent) {
	m.consumers[cName] = cChan
}

func (m *Manager) Deregister(cName string) {
	delete(m.consumers, cName)
}

// broadcastEvent takes an event and sends it to all registered consumers
func (m *Manager) broadcastEvent(e keyboard.ButtonEvent) {
	for _, c := range m.consumers {
		c <- e
	}
}

// Start ...
func (m *Manager) Start() {
	eChan := make(chan keyboard.ButtonEvent, 256)

	go func() {
		for evt := range eChan {
			m.broadcastEvent(evt)
		}
	}()

	k := NewKeylogger(&eChan)
	k.Start()

	// rk := NewReplayKeylogger(&eChan, nil)
	// rk.Start()
}
