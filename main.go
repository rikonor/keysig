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
	ttn.RegisterWith(k)

	top := metrics.NewTimeOfPress()
	top.RegisterWith(k)

	// Right now azul3d is blocking, therefore we shut off the logger with Ctrl+C
	// but also after a certain time we can write all our results to a csv file
	r := reports.New()
	go func() {
		// Wait 10 seconds before writing collected data to CSV
		time.Sleep(5 * time.Second)
		r.WriteToCSV(top)
	}()

	k.Start()
}
