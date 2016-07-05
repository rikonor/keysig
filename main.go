package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/rikonor/keysig/keylogger"
	"github.com/rikonor/keysig/metrics"
	"github.com/rikonor/keysig/reports"
)

func main() {
	m := keylogger.NewManager()
	r := reports.New()

	metrics.NewTimeToNext().RegisterWith(m).RegisterWithReporter(r)
	metrics.NewDurationOfPress().RegisterWith(m).RegisterWithReporter(r)
	metrics.NewKeyDistribution().RegisterWith(m).RegisterWithReporter(r)

	// Setup SIGTERM handler
	setTermHandler(r)

	r.TriggerPeriodically(3 * time.Minute)
	m.Start()
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
