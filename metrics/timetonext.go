package metrics

import (
	"fmt"
	"time"

	"azul3d.org/engine/keyboard"
	"github.com/rikonor/keysig/keylogger"
	"github.com/rikonor/keysig/reports"
	"github.com/rikonor/keysig/utils"
)

type timeToNextMetadata struct {
	// averageTime tracks the average transition time from key A to key B
	averageTime time.Duration
	// toNextCount tracks how many times you've transitioned from key A to key B
	toNextCount uint64
}

// timeToNextAggregate holds the average transition time between two keys
type timeToNextAggregate map[keyboard.Key]timeToNextMetadata

// TimeToNext metric type keeps track of the time delta between different characters
// E.g. the time from keyUp of key A and keyDown of key B
type TimeToNext struct {
	inputChan chan keyboard.ButtonEvent
	active    bool

	// timeToNextData contains all of our metrics data for key transitions
	timeToNextData map[keyboard.Key]timeToNextAggregate

	// currentlyPressed holds the currently pressed keys and when they were pressed
	currentlyPressed map[keyboard.Key]time.Time
}

func NewTimeToNext() *TimeToNext {
	return &TimeToNext{
		inputChan: make(chan keyboard.ButtonEvent),

		// Implementation specific data
		timeToNextData:   make(map[keyboard.Key]timeToNextAggregate),
		currentlyPressed: make(map[keyboard.Key]time.Time),
	}
}

func (m *TimeToNext) consumeStream() {
	go func() {
		for {
			time.Sleep(3 * time.Second)
			fmt.Println(m.currentlyPressed)
		}
	}()

	for evt := range m.inputChan {
		m.processEvent(evt)
	}
}

// RegisterWith registers with a keylogger
func (m *TimeToNext) RegisterWith(k *keylogger.Keylogger) {
	k.Register("timeToNext", m.inputChan)

	if !m.active {
		go m.consumeStream()
		m.active = true
	}
}

// RegisterWithReporter registers with a reporter
func (m *TimeToNext) RegisterWithReporter(r *reports.Reporter) {
	r.Register("timeToNext", m)
}

// Implementation

// handleDownEvent keep track of last time of Down event
func (m *TimeToNext) handleDownEvent(evt keyboard.ButtonEvent) {
	// Update currentlyPressed keys
	m.currentlyPressed[evt.Key] = time.Now()
}

// handleUpEvent keep track of last time of up event
// as well as update the average press time and press count
func (m *TimeToNext) handleUpEvent(evt keyboard.ButtonEvent) {
	// Since key was released, update the currently pressed keys
	delete(m.currentlyPressed, evt.Key)
}

func (m *TimeToNext) processEvent(evt keyboard.ButtonEvent) {
	switch evt.State {
	case keyboard.Down:
		m.handleDownEvent(evt)
	case keyboard.Up:
		m.handleUpEvent(evt)
	}
}

// Data collects our metrics data into a CSV compatible format
func (m *TimeToNext) Data() [][]string {
	return [][]string{
		{"first_name", "last_name", "username"},
		{"Poopie", "Pike", "rob"},
		{"Ken", "Thompson", "ken"},
		{"Robert", "Griesemer", "gri"},
	}
}

// WriteToCSV collects our metrics data and writes it to a CSV file
func (m *TimeToNext) WriteToCSV() {
	utils.WriteToCSV("timeToNext", m.Data())
}
