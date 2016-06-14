package metrics

import (
	"fmt"
	"time"

	"azul3d.org/engine/keyboard"
	"github.com/rikonor/keysig/keylogger"
	"github.com/rikonor/keysig/reports"
	"github.com/rikonor/keysig/utils"
)

type durationOfPressMetadata struct {
	averageTime time.Duration
	pressCount  uint64
}

// DurationOfPress metric type keeps track of the average time a key is pressed
// E.g. the time from keyUp of key A and keyDown of key A
type DurationOfPress struct {
	inputChan chan keyboard.ButtonEvent
	active    bool

	durationOfPressData map[keyboard.Key]durationOfPressMetadata

	lastDownTimes map[keyboard.Key]time.Time
	lastUpTimes   map[keyboard.Key]time.Time
}

func NewDurationOfPress() *DurationOfPress {
	return &DurationOfPress{
		inputChan: make(chan keyboard.ButtonEvent),

		// Implementation specific data
		durationOfPressData: make(map[keyboard.Key]durationOfPressMetadata),
		lastDownTimes:       make(map[keyboard.Key]time.Time),
		lastUpTimes:         make(map[keyboard.Key]time.Time),
	}
}

func (m *DurationOfPress) consumeStream() {
	for evt := range m.inputChan {
		m.processEvent(evt)
	}
}

// RegisterWith registers with a keylogger
func (m *DurationOfPress) RegisterWith(k *keylogger.Keylogger) *DurationOfPress {
	k.Register("durationOfPress", m.inputChan)

	if !m.active {
		go m.consumeStream()
		m.active = true
	}

	return m
}

// RegisterWithReporter registers with a reporter
func (m *DurationOfPress) RegisterWithReporter(r *reports.Reporter) {
	r.Register("durationOfPress", m)
}

// Implementation

func (m *DurationOfPress) String() string {
	output := ""

	for key, data := range m.durationOfPressData {
		output += fmt.Sprintf("%s %d %s\n", key, data.pressCount, data.averageTime)
	}

	return output + "\n\n"
}

// handleDownEvent keep track of last time of Down event
func (m *DurationOfPress) handleDownEvent(evt keyboard.ButtonEvent) {
	m.lastDownTimes[evt.Key] = time.Now()
}

// handleUpEvent keep track of last time of up event
// as well as update the average press time and press count
func (m *DurationOfPress) handleUpEvent(evt keyboard.ButtonEvent) {
	m.lastUpTimes[evt.Key] = time.Now()

	// Update the average
	currData := m.durationOfPressData[evt.Key]

	oldSum := currData.averageTime * time.Duration(currData.pressCount)
	newSum := oldSum + m.lastUpTimes[evt.Key].Sub(m.lastDownTimes[evt.Key])

	newPressCount := currData.pressCount + 1
	newAverage := time.Duration(uint64(newSum) / newPressCount)

	newPressData := durationOfPressMetadata{
		averageTime: newAverage,
		pressCount:  newPressCount,
	}

	m.durationOfPressData[evt.Key] = newPressData
}

func (m *DurationOfPress) processEvent(evt keyboard.ButtonEvent) {
	switch evt.State {
	case keyboard.Down:
		m.handleDownEvent(evt)
	case keyboard.Up:
		m.handleUpEvent(evt)
	}
}

// Data collects our metrics data into a CSV compatible format
func (m *DurationOfPress) Data() [][]string {
	// Convert durationOfPressData to [][]string for the reporter
	data := [][]string{
		{"key", "duration_of_press [ms]"},
	}

	for k, md := range m.durationOfPressData {
		l := []string{
			k.String(),
			utils.DurationToMSString(md.averageTime),
		}
		data = append(data, l)
	}

	return data
}

// WriteToCSV collects our metrics data and writes it to a CSV file
func (m *DurationOfPress) WriteToCSV() {
	utils.WriteToCSV("durationOfPress", m.Data())
}
