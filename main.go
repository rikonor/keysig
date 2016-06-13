package main

import (
	"time"

	"github.com/rikonor/keysig/keylogger"
	"github.com/rikonor/keysig/metrics"
	"github.com/rikonor/keysig/reports"
)

func main() {
	k := keylogger.New()

	ttn := metrics.NewTimeToNext()
	top := metrics.NewTimeOfPress()

	ms := []metrics.Metric{ttn, top}
	rs := []reports.Reportee{ttn, top}

	for _, m := range ms {
		m.RegisterWith(k)
	}

	// Right now azul3d is blocking, therefore we shut off the logger with Ctrl+C
	// but also after a certain time we can write all our results to a csv file
	r := reports.New()
	go func() {
		// Wait a few seconds before writing collected data to CSV
		time.Sleep(5 * time.Second)

		for _, p := range rs {
			r.WriteToCSV(p)
		}
	}()

	k.Start()
}
