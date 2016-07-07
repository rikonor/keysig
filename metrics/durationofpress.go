package metrics

import (
	"fmt"
	"time"

	"github.com/rikonor/keysig/keylogger"
	"github.com/rikonor/keysig/keylogger/keyboard"
	"github.com/rikonor/keysig/reports"
	"github.com/rikonor/keysig/utils"
)

const (
	// maxDuration is the keypress duration at which we assume
	// the keypress is not part of a normal typing flow (long pause, etc)
	maxDuration = time.Second

	// minPressCount is the minimum amount of key presses
	// that need to be recorded before a keypress data is assumed to be valid
	minPressCount uint64 = 2

	// maxDurationStd is the maximum allowed standard deviation
	// for keypress durations, anything above that means the values are too
	// spread out to have any statistical significance
	maxDurationStd = 50 * time.Millisecond
)

// DurationOfPress metric type keeps track of the average time a key is pressed
// E.g. the time from keyUp of key A and keyDown of key A
type DurationOfPress struct {
	inputChan chan keyboard.ButtonEvent
	active    bool

	durationOfPressData map[keyboard.Key]*utils.Stats

	lastDownTimes map[keyboard.Key]time.Time
	lastUpTimes   map[keyboard.Key]time.Time
}

func NewDurationOfPress() *DurationOfPress {
	return &DurationOfPress{
		inputChan: make(chan keyboard.ButtonEvent),

		// Implementation specific data
		durationOfPressData: make(map[keyboard.Key]*utils.Stats),
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
func (m *DurationOfPress) RegisterWith(k *keylogger.Manager) *DurationOfPress {
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
		output += fmt.Sprintf("%s %d %s\n", key, data.Count(), time.Duration(data.Mean()))
	}

	return output + "\n\n"
}

// handleDownEvent keep track of last time of Down event
func (m *DurationOfPress) handleDownEvent(evt keyboard.ButtonEvent) {
	m.lastDownTimes[evt.Key] = evt.Time()
}

// handleUpEvent keep track of last time of up event
// as well as update the average press time and press count
func (m *DurationOfPress) handleUpEvent(evt keyboard.ButtonEvent) {
	m.lastUpTimes[evt.Key] = evt.Time()

	pressDuration := m.lastUpTimes[evt.Key].Sub(m.lastDownTimes[evt.Key])
	if pressDuration > maxDuration {
		return
	}

	if m.durationOfPressData[evt.Key] == nil {
		m.durationOfPressData[evt.Key] = utils.NewStats()
	}

	m.durationOfPressData[evt.Key].Add(float64(pressDuration))
}

func (m *DurationOfPress) processEvent(evt keyboard.ButtonEvent) {
	switch evt.State {
	case keyboard.Down:
		m.handleDownEvent(evt)
	case keyboard.Up:
		m.handleUpEvent(evt)
	}
}

// DataForHeatMap outputs one line containing the average press durations
func (m *DurationOfPress) DataForHeatMap() [][]string {
	currLine := []string{}

	for _, k := range utils.OrderedKeys {
		md, ok := m.durationOfPressData[k]
		if !ok {
			// Should fill with zero instead of missing value
			currLine = append(currLine, "0.0")
			continue
		}

		// Check if keypress duration has been recorded enough times
		if md.Count() < minPressCount {
			currLine = append(currLine, "0.0")
			continue
		}

		// Require keypress durations to not be too spread out
		if md.Std() > float64(maxDurationStd) {
			currLine = append(currLine, "0.0")
			continue
		}

		currLine = append(currLine, utils.DurationToMSString(time.Duration(md.Mean())))
	}

	return [][]string{currLine}
}

// Data collects our metrics data into a CSV compatible format
// DEPRECATED: Currently using DataForHeatMap instead
func (m *DurationOfPress) Data() [][]string {
	// Convert durationOfPressData to [][]string for the reporter
	data := [][]string{
		{"key", "duration_of_press [ms]"},
	}

	// Iterating over OrderedKeys so our data is sorted
	// Notice OrderedKeys doesn't contain every possible key (only letters)
	for _, k := range utils.OrderedKeys {
		// check if data was captured for current key
		if md, ok := m.durationOfPressData[k]; ok {
			l := []string{
				k.String(),
				utils.DurationToMSString(time.Duration(md.Mean())),
			}
			data = append(data, l)
		}
	}

	return data
}

// WriteToCSV collects our metrics data and writes it to a CSV file
func (m *DurationOfPress) WriteToCSV() {
	utils.WriteToCSV("durationOfPress", m.DataForHeatMap())
}
