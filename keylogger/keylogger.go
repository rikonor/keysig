package keylogger

import "github.com/rikonor/keysig/keylogger/keyboard"

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

// Start ...
func (k *Keylogger) Start() {
	eChan := make(chan keyboard.ButtonEvent, 256)

	go func() {
		for evt := range eChan {
			k.broadcastEvent(evt)
		}
	}()

	// osxk := osxkeylogger.NewOSXKeylogger(&eChan)
	// osxk.Start()

	// a3dk := azul3dkeylogger.NewAzul3DKeylogger(&eChan)
	// a3dk.Start()

	rk := NewReplayKeylogger(&eChan, nil)
	rk.Start()
}
