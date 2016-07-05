package keylogger

import (
	"fmt"
	"time"

	"github.com/MarinX/keylogger"
	"github.com/rikonor/keysig/keylogger/keyboard"
)

// Linux keylogger

type Keylogger struct {
	outputChannel *chan keyboard.ButtonEvent
}

func NewKeylogger(outputChannel *chan keyboard.ButtonEvent) *Keylogger {
	return &Keylogger{outputChannel: outputChannel}
}

func (k *Keylogger) handleButtonEvent(keyCode int, stateCode int) {
	kc := convertKeyCode(keyCode)
	// Skip invalid keys
	if kc == keyboard.Invalid {
		return
	}

	s := convertStateCode(stateCode)
	// Skip invalid state
	if s == keyboard.InvalidState {
		return
	}

	evt := keyboard.ButtonEvent{
		T:     time.Now(),
		Key:   kc,
		State: s,
	}

	*k.outputChannel <- evt
}

func (k *Keylogger) Start() {
	devs, err := keylogger.NewDevices()
	if err != nil {
		fmt.Println(err)
		return
	}

	// for _, val := range devs {
	// 	fmt.Println("Id->", val.Id, "Device->", val.Name)
	// }

	rd := keylogger.NewKeyLogger(devs[2])

	eventStream, err := rd.Read()
	if err != nil {
		fmt.Println(err)
		return
	}

	for evt := range eventStream {
		// Filter non-keyboard events
		if evt.Type != keylogger.EV_KEY {
			continue
		}

		k.handleButtonEvent(int(evt.Code), int(evt.Value))
	}
}
