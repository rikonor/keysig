package main

import (
	"time"

	"github.com/rikonor/keysig/keylogger"
	"github.com/rikonor/keysig/metrics"
	"github.com/rikonor/keysig/reports"
)

func main() {
	k := keylogger.New()
	r := reports.New()

	ttn := metrics.NewTimeToNext()
	ttn.RegisterWith(k)
	ttn.RegisterWithReporter(r)

	top := metrics.NewTimeOfPress()
	top.RegisterWith(k)
	top.RegisterWithReporter(r)

	// Right now azul3d is blocking, therefore we shut off the logger with Ctrl+C
	// but also after a certain time we can write all our results to a csv file
	go func() {
		// Wait a few seconds before writing collected data to CSV
		time.Sleep(5 * time.Second)
		r.CollectReports()
	}()

	k.Start()
}
