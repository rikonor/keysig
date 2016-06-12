package main

import (
	"github.com/rikonor/keysig/keylogger"
	"github.com/rikonor/keysig/metrics"
)

func main() {
	k := keylogger.New()

	metrics.NewCharDiff().RegisterWith(k)

	// // Register all of your typing statistics
	// evts := make(chan keyboard.ButtonEvent)
	// go func() {
	// 	for evt := range evts {
	// 		fmt.Println("Got something", evt)
	// 	}
	// }()
	//
	// k.Register("myChan", evts)

	k.Start()
}
