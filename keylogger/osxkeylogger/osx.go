package osxkeylogger

import (
	"fmt"

	"github.com/rikonor/keysig/keylogger/keyboard"
)

// Based on github.com/caseyscarborough/keylogger

/*
#cgo LDFLAGS: -framework ApplicationServices -framework Carbon
#include "keylogger.h"

extern void handleButtonEvent(int k);

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

		handleButtonEvent((int)keyCode);

		return event;
}
*/
import "C"

var globalOutputChannel *chan keyboard.ButtonEvent

//export handleButtonEvent
func handleButtonEvent(keyCode C.int) {
	// TODO
	// Convert from the int we get to an actual keyboard.Key
	// Dedup still pressed keys. If a key is pressed for more then 1s then it will continue being reported
	// Find the current key state based on previous events (the first event should always mean pressing down)

	k := convertKeyCode(int(keyCode))
	fmt.Println("Got", k)

	evt := keyboard.ButtonEvent{}

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

// Utils

var keyCodeConversionTable = map[int]keyboard.Key{
	0x00: keyboard.A,
	0x01: keyboard.S,
	0x02: keyboard.D,
	0x03: keyboard.F,
	0x04: keyboard.H,
	0x05: keyboard.G,
	0x06: keyboard.Z,
	0x07: keyboard.X,
	0x08: keyboard.C,
	0x09: keyboard.V,
	0x0B: keyboard.B,
	0x0C: keyboard.Q,
	0x0D: keyboard.W,
	0x0E: keyboard.E,
	0x0F: keyboard.R,
	0x10: keyboard.Y,
	0x11: keyboard.T,
	0x12: keyboard.One,
	0x13: keyboard.Two,
	0x14: keyboard.Three,
	0x15: keyboard.Four,
	0x16: keyboard.Six,
	0x17: keyboard.Five,
	0x18: keyboard.Equals,
	0x19: keyboard.Nine,
	0x1A: keyboard.Seven,
	0x1B: keyboard.Invalid, // Minus
	0x1C: keyboard.Eight,
	0x1D: keyboard.Zero,
	0x1E: keyboard.RightBracket,
	0x1F: keyboard.O,
	0x20: keyboard.U,
	0x21: keyboard.LeftBracket,
	0x22: keyboard.I,
	0x23: keyboard.P,
	0x25: keyboard.L,
	0x26: keyboard.J,
	0x27: keyboard.Invalid, // Quote
	0x28: keyboard.K,
	0x29: keyboard.Semicolon,
	0x2A: keyboard.BackSlash,
	0x2B: keyboard.Comma,
	0x2C: keyboard.Invalid, // Slash
	0x2D: keyboard.N,
	0x2E: keyboard.M,
	0x2F: keyboard.Period,
	0x32: keyboard.Invalid, // Grave
	0x41: keyboard.Invalid, // KeypadDecimal
	0x43: keyboard.Invalid, // KeypadMultiply
	0x45: keyboard.Invalid, // KeypadPlus
	0x47: keyboard.Invalid, // KeypadClear
	0x4B: keyboard.Invalid, // KeypadDivide
	0x4C: keyboard.Invalid, // KeypadEnter
	0x4E: keyboard.Invalid, // KeypadMinus
	0x51: keyboard.Invalid, // KeypadEquals
	0x52: keyboard.Invalid, // Keypad0
	0x53: keyboard.Invalid, // Keypad1
	0x54: keyboard.Invalid, // Keypad2
	0x55: keyboard.Invalid, // Keypad3
	0x56: keyboard.Invalid, // Keypad4
	0x57: keyboard.Invalid, // Keypad5
	0x58: keyboard.Invalid, // Keypad6
	0x59: keyboard.Invalid, // Keypad7
	0x5B: keyboard.Invalid, // Keypad8
	0x5C: keyboard.Invalid, // Keypad9
	0x24: keyboard.Enter,
	0x30: keyboard.Tab,
	0x31: keyboard.Space,
	0x33: keyboard.Delete,
	0x35: keyboard.Escape,
	0x37: keyboard.Invalid, // Command
	0x38: keyboard.Invalid, // Shift
	0x39: keyboard.CapsLock,
	0x3A: keyboard.Invalid, // Option
	0x3B: keyboard.Invalid, // Control
	0x3C: keyboard.RightShift,
	0x3D: keyboard.Invalid, // RightOption
	0x3E: keyboard.RightCtrl,
	0x3F: keyboard.Invalid, // Function
	0x40: keyboard.F17,
	0x48: keyboard.Invalid, // VolumeUp
	0x49: keyboard.Invalid, // VolumeDown
	0x4A: keyboard.Invalid, // Mute
	0x4F: keyboard.F18,
	0x50: keyboard.F19,
	0x5A: keyboard.F20,
	0x60: keyboard.F5,
	0x61: keyboard.F6,
	0x62: keyboard.F7,
	0x63: keyboard.F3,
	0x64: keyboard.F8,
	0x65: keyboard.F9,
	0x67: keyboard.F11,
	0x69: keyboard.F13,
	0x6A: keyboard.F16,
	0x6B: keyboard.F14,
	0x6D: keyboard.F10,
	0x6F: keyboard.F12,
	0x71: keyboard.F15,
	0x72: keyboard.Help,
	0x73: keyboard.Home,
	0x74: keyboard.PageUp,
	0x75: keyboard.Invalid, // ForwardDelete
	0x76: keyboard.F4,
	0x77: keyboard.End,
	0x78: keyboard.F2,
	0x79: keyboard.PageDown,
	0x7A: keyboard.F1,
	0x7B: keyboard.ArrowLeft,
	0x7C: keyboard.ArrowRight,
	0x7D: keyboard.ArrowDown,
	0x7E: keyboard.ArrowUp,
}

func convertKeyCode(keyCode int) keyboard.Key {
	// Search for keyCode in the conversion table
	k, ok := keyCodeConversionTable[keyCode]
	if !ok {
		return keyboard.Invalid
	}
	return k
}
