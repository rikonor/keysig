package linuxkeylogger

import (
	"fmt"
	"log"
	"strings"
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

func (k *Keylogger) setupDevice() *keylogger.KeyLogger {
	devs, err := keylogger.NewDevices()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// Search for the keyboard device
	devIdx := -1
	for idx, val := range devs {
		if strings.Contains(val.Name, "keyboard") {
			devIdx = idx
			break
		}
	}
	if devIdx == -1 {
		return nil
	}

	return keylogger.NewKeyLogger(devs[devIdx])
}

func (k *Keylogger) Start() {
	rd := k.setupDevice()
	if rd == nil {
		log.Fatalln("Failed to find keyboard device. Aborting..")
	}

	eventStream, err := rd.Read()
	if err != nil {
		log.Fatalln(err)
	}

	for evt := range eventStream {
		// Filter non-keyboard events
		if evt.Type != keylogger.EV_KEY {
			continue
		}

		k.handleButtonEvent(int(evt.Code), int(evt.Value))
	}
}
