package metrics

import (
	"time"

	"github.com/rikonor/keysig/keylogger"
	"github.com/rikonor/keysig/keylogger/keyboard"
	"github.com/rikonor/keysig/reports"
	"github.com/rikonor/keysig/utils"
)

const (
	// maxTransitionDuration is the transition duration at which we assume
	// the transition is not part of a normal typing flow (long pause, etc)
	maxTransitionDuration = time.Second

	// minTransitionsCount is the minimum amount of transitions
	// that need to be recorded before a transition is assumed to be valid
	minTransitionsCount uint64 = 2

	// maxTransitionStd is the maximum allowed standard deviation
	// for transition durations, anything above that means the values are too
	// spread out to have any statistical significance
	maxTransitionStd = 50 * time.Millisecond
)

type PressMetadata struct {
	key  keyboard.Key
	time time.Time
}

// timeToNextAggregate holds the average transition time between two keys
type timeToNextAggregate map[keyboard.Key]*utils.Stats

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
func (m *TimeToNext) RegisterWith(k *keylogger.Manager) *TimeToNext {
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

	// Ignore very long transitions, they might represent long pauses in typing
	transitionDuration := evt.Time().Sub(m.lastUpEvent.time)
	if transitionDuration > maxTransitionDuration {
		return
	}

	// If this is the first time this transition has occured
	// initialize the data structure for it
	if _, ok := m.timeToNextData[m.lastUpEvent.key]; !ok {
		m.timeToNextData[m.lastUpEvent.key] = make(timeToNextAggregate)
	}

	if m.timeToNextData[m.lastUpEvent.key][evt.Key] == nil {
		m.timeToNextData[m.lastUpEvent.key][evt.Key] = utils.NewStats()
	}

	m.timeToNextData[m.lastUpEvent.key][evt.Key].Add(float64(transitionDuration))
}

// handleUpEvent keep track of last time of up event
func (m *TimeToNext) handleUpEvent(evt keyboard.ButtonEvent) {
	// Since key was released, update the currently pressed keys
	m.lastUpEvent = PressMetadata{
		key:  evt.Key,
		time: evt.Time(),
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

func (m *TimeToNext) Data() [][]string {
	data := [][]string{}

	numOfKeys := len(utils.OrderedKeys)

	for _, fromKey := range utils.OrderedKeys {
		currLine := []string{}

		// Check if fromKey has been previously recorded
		_, ok := m.timeToNextData[fromKey]
		if !ok {
			// Shoud fill row with zeros
			for i := 0; i < numOfKeys; i++ {
				currLine = append(currLine, "0.0")
			}
			data = append(data, currLine)
			// Continue to next key
			continue
		}

		for _, toKey := range utils.OrderedKeys {
			// Check if toKey has been previoulsy recorded for fromKey
			md, ok := m.timeToNextData[fromKey][toKey]
			if !ok {
				// Should fill with a zero instead of missing value
				currLine = append(currLine, "0.0")
				continue
			}

			// Check if transition has been recorded enough times
			if md.Count() < minTransitionsCount {
				currLine = append(currLine, "0.0")
				continue
			}

			// Require transition times to not be too spread out
			if md.Std() > float64(maxTransitionStd) {
				currLine = append(currLine, "0.0")
				continue
			}

			currLine = append(currLine, utils.DurationToMSString(time.Duration(md.Mean())))
		}

		data = append(data, currLine)
	}

	// Add headers on top and left [the key names]
	// data = append(utils.OrderedKeys, data...)

	return data
}

// WriteToCSV collects our metrics data and writes it to a CSV file
func (m *TimeToNext) WriteToCSV() {
	utils.WriteToCSV("timeToNext", m.Data())
}
