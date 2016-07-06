Keylogger
---

This package provides a cross-platform keylogger.
It's a wrapper around the following:

* Linux - [MarinX/keylogger](https://github.com/MarinX/keylogger)
* OSX - [caseyscarborough/keylogger](https://github.com/caseyscarborough/keylogger)
* Windows - [azul3d/engine](https://github.com/azul3d/engine) [To be replaced]

Given a channel, it will stream keystroke events based on the interface defined in [Azul3D](github.com/rikonor/keysig/blob/master/keylogger/keyboard/events.go), where each event contains the `Key`, `State` (press/release) and `Time`.

Requires `root`.

##### Example

```
import (
	"fmt"

	"github.com/rikonor/keysig/keylogger"
	"github.com/rikonor/keysig/keylogger/keyboard"
)

func main() {
	m := keylogger.NewManager()

	ch := make(chan keyboard.ButtonEvent)
	m.Register("example", ch)

	go func() {
		for evt := range ch {
			fmt.Println(evt)
		}
	}()

	m.Start()
}
```

##### Output

```
ButtonEvent(Key=H, State=Down, Time=2016-07-06 16:45:45.281751354 -0400 EDT)
ButtonEvent(Key=H, State=Up, Time=2016-07-06 16:45:45.350629246 -0400 EDT)
ButtonEvent(Key=E, State=Down, Time=2016-07-06 16:45:45.443853526 -0400 EDT)
ButtonEvent(Key=E, State=Up, Time=2016-07-06 16:45:45.511998649 -0400 EDT)
```
