package metrics

import (
	"fmt"

	"github.com/rikonor/keysig/keylogger"
	"github.com/rikonor/keysig/keylogger/keyboard"
	"github.com/rikonor/keysig/reports"
	"github.com/rikonor/keysig/utils"
)

type keyDistributionMetadata struct {
	pressCount uint64
}

// KeyDistribution metric type keeps track of the average time a key is pressed
// E.g. the time from keyUp of key A and keyDown of key A
type KeyDistribution struct {
	inputChan chan keyboard.ButtonEvent
	active    bool

	keyDistributionData map[keyboard.Key]keyDistributionMetadata
	totalPressCount     uint64
}

func NewKeyDistribution() *KeyDistribution {
	return &KeyDistribution{
		inputChan: make(chan keyboard.ButtonEvent),

		// Implementation specific data
		keyDistributionData: make(map[keyboard.Key]keyDistributionMetadata),
	}
}

func (m *KeyDistribution) consumeStream() {
	for evt := range m.inputChan {
		m.processEvent(evt)
	}
}

// RegisterWith registers with a keylogger
func (m *KeyDistribution) RegisterWith(k *keylogger.Manager) *KeyDistribution {
	k.Register("keyDistribution", m.inputChan)

	if !m.active {
		go m.consumeStream()
		m.active = true
	}

	return m
}

// RegisterWithReporter registers with a reporter
func (m *KeyDistribution) RegisterWithReporter(r *reports.Reporter) {
	r.Register("keyDistribution", m)
}

// Implementation

func (m *KeyDistribution) String() string {
	output := ""

	output += fmt.Sprintf("Total=%d", m.totalPressCount)

	for key, data := range m.keyDistributionData {
		output += fmt.Sprintf("%s %d\n", key, data.pressCount)
	}

	return output + "\n\n"
}

// handleDownEvent keep track of last time of Down event
func (m *KeyDistribution) handleDownEvent(evt keyboard.ButtonEvent) {
	m.totalPressCount++

	currData := m.keyDistributionData[evt.Key]
	m.keyDistributionData[evt.Key] = keyDistributionMetadata{
		pressCount: currData.pressCount + 1,
	}
}

func (m *KeyDistribution) processEvent(evt keyboard.ButtonEvent) {
	switch evt.State {
	case keyboard.Down:
		m.handleDownEvent(evt)
	}
}

// Data collects our metrics data into a CSV compatible format
func (m *KeyDistribution) Data() [][]string {
	// Convert keyDistribution to [][]string for the reporter
	data := [][]string{
		{"key", "distribution"},
	}

	for k, md := range m.keyDistributionData {
		// Calculate the actual distribution
		keyDistribution := float64(md.pressCount) / float64(m.totalPressCount)

		l := []string{
			k.String(),
			fmt.Sprint(keyDistribution),
		}
		data = append(data, l)
	}

	return data
}

// WriteToCSV collects our metrics data and writes it to a CSV file
func (m *KeyDistribution) WriteToCSV() {
	utils.WriteToCSV("keyDistribution", m.Data())
}
