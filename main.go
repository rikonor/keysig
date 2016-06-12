package main

import (
	"github.com/rikonor/keysig/keylogger"
	"github.com/rikonor/keysig/metrics"
)

func main() {
	k := keylogger.New()

	metrics.NewTimeToNext().RegisterWith(k)
	metrics.NewTimeOfPress().RegisterWith(k)

	k.Start()
}
