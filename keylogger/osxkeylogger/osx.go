package osxkeylogger

import (
	"fmt"
	"time"

	"github.com/rikonor/keysig/keylogger/keyboard"
)

// Based on github.com/caseyscarborough/keylogger

/*
#cgo LDFLAGS: -framework ApplicationServices -framework Carbon
#include "keylogger.h"

typedef enum State { Invalid, Down, Up } State;

extern void handleButtonEvent(int k, State s);

static inline void start_logger() {
    // Create an event tap to retrieve keypresses.
    CGEventMask eventMask = (
        CGEventMaskBit(kCGEventKeyDown) |
        CGEventMaskBit(kCGEventKeyUp) |
        CGEventMaskBit(kCGEventFlagsChanged)
    );
    CFMachPortRef eventTap = CGEventTapCreate(
        kCGSessionEventTap, kCGHeadInsertEventTap, 0, eventMask, CGEventCallback, NULL
    );

    // Exit the program if unable to create the event tap.
    if(!eventTap) {
        fprintf(stderr, "ERROR: Unable to create event tap.\n");
        exit(1);
    }

    // Create a run loop source and add enable the event tap.
    CFRunLoopSourceRef runLoopSource = CFMachPortCreateRunLoopSource(kCFAllocatorDefault, eventTap, 0);
    CFRunLoopAddSource(CFRunLoopGetCurrent(), runLoopSource, kCFRunLoopCommonModes);
    CGEventTapEnable(eventTap, true);

    // Get the current time and open the logfile.
    time_t result = time(NULL);
    logfile = fopen(logfileLocation, "a");

    if (!logfile) {
        fprintf(stderr, "ERROR: Unable to open log file. Ensure that you have the proper permissions.\n");
        exit(1);
    }

    CFRunLoopRun();
}

// The following callback method is invoked on every keypress.
static inline CGEventRef CGEventCallback(CGEventTapProxy proxy, CGEventType type, CGEventRef event, void *refcon) {
    if (type != kCGEventKeyDown && type != kCGEventFlagsChanged && type != kCGEventKeyUp) { return event; }

    // Retrieve the incoming keycode.
    CGKeyCode keyCode = (CGKeyCode) CGEventGetIntegerValueField(event, kCGKeyboardEventKeycode);

		State s;
		if (type == kCGEventKeyDown) s = Down;
		if (type == kCGEventKeyUp) s = Up;
		if (type == kCGEventFlagsChanged) s = Invalid;

		handleButtonEvent((int)keyCode, s);

		return event;
}
*/
import "C"

var globalOutputChannel *chan keyboard.ButtonEvent

//export handleButtonEvent
func handleButtonEvent(keyCode C.int, s C.State) {
	// TODO
	// [Don] Convert from the int we get to an actual keyboard.Key
	// Dedup still pressed keys. If a key is pressed for more then 1s then it will continue being reported
	// Find the current key state based on previous events (the first event should always mean pressing down)

	k := convertKeyCode(int(keyCode))
	fmt.Println("Got", k, s)

	evt := keyboard.ButtonEvent{
		T:     time.Now(),
		Key:   k,
		State: keyboard.Up,
	}

	*globalOutputChannel <- evt
}

type OSXKeylogger struct {
	outputChannel *chan keyboard.ButtonEvent
}

func NewOSXKeylogger(outputChannel *chan keyboard.ButtonEvent) *OSXKeylogger {
	// Assign the global outputChannel
	globalOutputChannel = outputChannel
	return &OSXKeylogger{outputChannel: outputChannel}
}

func (k *OSXKeylogger) Start() {
	C.start_logger()
}
