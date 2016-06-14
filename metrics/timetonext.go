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
	transitionsCount uint64
}

type PressMetadata struct {
	key  keyboard.Key
	time time.Time
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

	lastUpEvent PressMetadata
}

func NewTimeToNext() *TimeToNext {
	return &TimeToNext{
		inputChan: make(chan keyboard.ButtonEvent),

		// Implementation specific data
		timeToNextData: make(map[keyboard.Key]timeToNextAggregate),
	}
}

func (m *TimeToNext) consumeStream() {
	for evt := range m.inputChan {
		m.processEvent(evt)
	}
}

// RegisterWith registers with a keylogger
func (m *TimeToNext) RegisterWith(k *keylogger.Keylogger) *TimeToNext {
	k.Register("timeToNext", m.inputChan)

	if !m.active {
		go m.consumeStream()
		m.active = true
	}

	return m
}

// RegisterWithReporter registers with a reporter
func (m *TimeToNext) RegisterWithReporter(r *reports.Reporter) {
	r.Register("timeToNext", m)
}

// Implementation

// handleDownEvent keep track of last time of Down event
func (m *TimeToNext) handleDownEvent(evt keyboard.ButtonEvent) {
	// Skip the first time because lastUpEvent won't contain a valid value
	if m.lastUpEvent == (PressMetadata{}) {
		return
	}

	// If this is the first time this transition has occured
	// initialize the data structure for it
	_, ok := m.timeToNextData[m.lastUpEvent.key][evt.Key]
	if !ok {
		m.timeToNextData[m.lastUpEvent.key] = make(timeToNextAggregate)
	}

	// Update the values transitionsCount and averageTime for this transition
	oldTransitionsCount := m.timeToNextData[m.lastUpEvent.key][evt.Key].transitionsCount

	ttnm := timeToNextMetadata{
		averageTime:      time.Now().Sub(m.lastUpEvent.time),
		transitionsCount: oldTransitionsCount + 1,
	}

	m.timeToNextData[m.lastUpEvent.key][evt.Key] = ttnm

	fmt.Println(
		m.lastUpEvent.key,
		"->",
		evt.Key,
		fmt.Sprintf("%dms\t", ttnm.averageTime.Nanoseconds()/(1000*1000)),
		fmt.Sprintf("[%d times]", m.timeToNextData[m.lastUpEvent.key][evt.Key].transitionsCount),
	)
}

// handleUpEvent keep track of last time of up event
func (m *TimeToNext) handleUpEvent(evt keyboard.ButtonEvent) {
	// Since key was released, update the currently pressed keys
	m.lastUpEvent = PressMetadata{
		key:  evt.Key,
		time: time.Now(),
	}
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
