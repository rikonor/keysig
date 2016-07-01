package main

import (
	"os"
	"os/signal"

	"github.com/rikonor/keysig/keylogger"
	"github.com/rikonor/keysig/metrics"
	"github.com/rikonor/keysig/reports"
)

func main() {
	m := keylogger.NewManager()
	r := reports.New()

	metrics.NewTimeToNext().RegisterWith(m).RegisterWithReporter(r)
	metrics.NewDurationOfPress().RegisterWith(m).RegisterWithReporter(r)

	// Right now azul3d is blocking, therefore we shut off the logger with Ctrl+C
	setTermHandler(r)

	m.Start()

	// We can also shut off the logger by closing the window
	r.CollectReports()
}

func setTermHandler(r *reports.Reporter) {
	go func() {
		// Wait for SIGTERM and collect the reports
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c

		r.CollectReports()
		os.Exit(0)
	}()
}
