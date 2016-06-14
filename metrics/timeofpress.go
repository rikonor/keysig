package metrics

import (
	"fmt"
	"time"

	"azul3d.org/engine/keyboard"
	"github.com/rikonor/keysig/keylogger"
	"github.com/rikonor/keysig/reports"
	"github.com/rikonor/keysig/utils"
)

type timeOfPressMetadata struct {
	averageTime time.Duration
	pressCount  uint64
}

// TimeOfPress metric type keeps track of the average time a key is pressed
// E.g. the time from keyUp of key A and keyDown of key A
type TimeOfPress struct {
	inputChan chan keyboard.ButtonEvent
	active    bool

	timeOfPressData map[keyboard.Key]timeOfPressMetadata

	lastDownTimes map[keyboard.Key]time.Time
	lastUpTimes   map[keyboard.Key]time.Time
}

func NewTimeOfPress() *TimeOfPress {
	return &TimeOfPress{
		inputChan: make(chan keyboard.ButtonEvent),

		// Implementation specific data
		timeOfPressData: make(map[keyboard.Key]timeOfPressMetadata),
		lastDownTimes:   make(map[keyboard.Key]time.Time),
		lastUpTimes:     make(map[keyboard.Key]time.Time),
	}
}

func (m *TimeOfPress) consumeStream() {
	for evt := range m.inputChan {
		m.processEvent(evt)
	}
}

// RegisterWith registers with a keylogger
func (m *TimeOfPress) RegisterWith(k *keylogger.Keylogger) *TimeOfPress {
	k.Register("timeOfPress", m.inputChan)

	if !m.active {
		go m.consumeStream()
		m.active = true
	}

	return m
}

// RegisterWithReporter registers with a reporter
func (m *TimeOfPress) RegisterWithReporter(r *reports.Reporter) {
	r.Register("timeOfPress", m)
}

// Implementation

func (m *TimeOfPress) String() string {
	output := ""

	for key, data := range m.timeOfPressData {
		output += fmt.Sprintf("%s %d %s\n", key, data.pressCount, data.averageTime)
	}

	return output + "\n\n"
}

// handleDownEvent keep track of last time of Down event
func (m *TimeOfPress) handleDownEvent(evt keyboard.ButtonEvent) {
	m.lastDownTimes[evt.Key] = time.Now()
}

// handleUpEvent keep track of last time of up event
// as well as update the average press time and press count
func (m *TimeOfPress) handleUpEvent(evt keyboard.ButtonEvent) {
	m.lastUpTimes[evt.Key] = time.Now()

	// Update the average
	currData := m.timeOfPressData[evt.Key]

	oldSum := currData.averageTime * time.Duration(currData.pressCount)
	newSum := oldSum + m.lastUpTimes[evt.Key].Sub(m.lastDownTimes[evt.Key])

	newPressCount := currData.pressCount + 1
	newAverage := time.Duration(uint64(newSum) / newPressCount)

	newPressData := timeOfPressMetadata{
		averageTime: newAverage,
		pressCount:  newPressCount,
	}

	m.timeOfPressData[evt.Key] = newPressData
}

func (m *TimeOfPress) processEvent(evt keyboard.ButtonEvent) {
	switch evt.State {
	case keyboard.Down:
		m.handleDownEvent(evt)
	case keyboard.Up:
		m.handleUpEvent(evt)
	}
}

// Data collects our metrics data into a CSV compatible format
func (m *TimeOfPress) Data() [][]string {
	// Convert timeOfPressData to [][]string for the reporter
	data := [][]string{
		{"key", "time_of_press"},
	}

	for k, md := range m.timeOfPressData {
		l := []string{
			k.String(),
			utils.DurationToMSString(md.averageTime),
		}
		data = append(data, l)
	}

	return data
}

// WriteToCSV collects our metrics data and writes it to a CSV file
func (m *TimeOfPress) WriteToCSV() {
	utils.WriteToCSV("timeOfPress", m.Data())
}
