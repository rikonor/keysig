package linuxkeylogger

import "github.com/rikonor/keysig/keylogger/keyboard"

var keyCodeConversionTable = map[int]keyboard.Key{
	0x01: keyboard.Escape,
	0x02: keyboard.One,
	0x03: keyboard.Two,
	0x04: keyboard.Three,
	0x05: keyboard.Four,
	0x06: keyboard.Five,
	0x07: keyboard.Six,
	0x08: keyboard.Seven,
	0x09: keyboard.Eight,
	0x0A: keyboard.Nine,
	0x0B: keyboard.Zero,
	0x0C: keyboard.Dash,
	0x0D: keyboard.Equals,
	0x0E: keyboard.Backspace,
	0x0F: keyboard.Tab,
	0x10: keyboard.Q,
	0x11: keyboard.W,
	0x12: keyboard.E,
	0x13: keyboard.R,
	0x14: keyboard.T,
	0x15: keyboard.Y,
	0x16: keyboard.U,
	0x17: keyboard.I,
	0x18: keyboard.O,
	0x19: keyboard.P,
	0x1A: keyboard.LeftBracket,
	0x1B: keyboard.RightBracket,
	0x1C: keyboard.Enter,
	0x1D: keyboard.LeftCtrl,
	0x1E: keyboard.A,
	0x1F: keyboard.S,
	0x20: keyboard.D,
	0x21: keyboard.F,
	0x22: keyboard.G,
	0x23: keyboard.H,
	0x24: keyboard.J,
	0x25: keyboard.K,
	0x26: keyboard.L,
	0x27: keyboard.Semicolon,
	0x28: keyboard.Apostrophe,
	0x29: keyboard.Tilde,
	0x2A: keyboard.LeftShift,
	0x2B: keyboard.ForwardSlash,
	0x2C: keyboard.Z,
	0x2D: keyboard.X,
	0x2E: keyboard.C,
	0x2F: keyboard.V,
	0x30: keyboard.B,
	0x31: keyboard.N,
	0x32: keyboard.M,
	0x33: keyboard.Comma,
	0x34: keyboard.Period,
	0x35: keyboard.BackSlash,
	0x36: keyboard.RightShift,
	0x38: keyboard.LeftAlt,
	0x39: keyboard.Space,
	0x3A: keyboard.CapsLock,
	0x3B: keyboard.F1,
	0x3C: keyboard.F2,
	0x3D: keyboard.F3,
	0x3E: keyboard.F4,
	0x3F: keyboard.F5,
	0x40: keyboard.F6,
	0x41: keyboard.F7,
	0x42: keyboard.F8,
	0x43: keyboard.F9,
	0x44: keyboard.F10,
	0x45: keyboard.NumLock,
	0x46: keyboard.ScrollLock,
	0x47: keyboard.Home,
	0x57: keyboard.F11,
	0x58: keyboard.F12,
	0x64: keyboard.RightAlt,
	0x67: keyboard.ArrowUp,
	0x69: keyboard.ArrowLeft,
	0x6A: keyboard.ArrowRight,
	0x6C: keyboard.ArrowDown,
	0x6F: keyboard.Delete,
}

func convertKeyCode(keyCode int) keyboard.Key {
	// Search for keyCode in the conversion table
	k, ok := keyCodeConversionTable[keyCode]
	if !ok {
		return keyboard.Invalid
	}
	return k
}

var stateCodeConversionTable = map[int]keyboard.State{
	0: keyboard.Up,
	1: keyboard.Down,
	2: keyboard.InvalidState,
}

func convertStateCode(stateCode int) keyboard.State {
	// Search for stateCode in the conversion table
	s, ok := stateCodeConversionTable[stateCode]
	if !ok {
		return keyboard.InvalidState
	}
	return s
}
