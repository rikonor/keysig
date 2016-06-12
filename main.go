package main

import (
	"fmt"

	"azul3d.org/engine/keyboard"
	"github.com/rikonor/keysig/keylogger"
)

func main() {
	k := keylogger.New()

	// Register all of your typing statistics
	evts := make(chan keyboard.ButtonEvent)
	go func() {
		for evt := range evts {
			fmt.Println("Got something", evt)
		}
	}()

	k.Register("myChan", evts)

	k.Start()
}
