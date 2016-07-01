package keylogger

import (
	"math/rand"
	"time"

	"github.com/rikonor/keysig/keylogger/keyboard"
	"github.com/rikonor/keysig/utils"
)

// ReplayKeylogger can replay a sequence of button events
// If not provided with one it will use a randomly generated sequence
type ReplayKeylogger struct {
	outputChannel       *chan keyboard.ButtonEvent
	buttonEventSequence []keyboard.ButtonEvent
}

func NewReplayKeylogger(oc *chan keyboard.ButtonEvent, bes []keyboard.ButtonEvent) *ReplayKeylogger {
	// fallback to random sequence
	if bes == nil {
		bes = randomButtonEventSequence()
	}

	return &ReplayKeylogger{
		outputChannel:       oc,
		buttonEventSequence: bes,
	}
}

func (k *ReplayKeylogger) Start() {
	for _, be := range k.buttonEventSequence {
		time.Sleep(time.Millisecond)
		*k.outputChannel <- be
	}
}

// Utils

var sequenceLength = 1000

func randomButtonEventSequence() []keyboard.ButtonEvent {
	seq := []keyboard.ButtonEvent{}
	for i := 0; i < sequenceLength; i++ {
		seq = append(seq, randomButtonEvent())
	}
	return seq
}

func randomButtonEvent() keyboard.ButtonEvent {
	randomKeyIdx := rand.Intn(len(utils.OrderedKeys))
	randomKey := utils.OrderedKeys[randomKeyIdx]

	randomStateIdx := rand.Intn(2)
	randomState := []keyboard.State{keyboard.Down, keyboard.Up}[randomStateIdx]

	return keyboard.ButtonEvent{
		T:     time.Now(),
		Key:   randomKey,
		State: randomState,
	}
}
