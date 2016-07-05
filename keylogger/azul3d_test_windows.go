package keylogger

import (
	"reflect"
	"testing"
	"time"

	azul3d "azul3d.org/engine/keyboard"
	"github.com/rikonor/keysig/keylogger/keyboard"
)

func TestCanConvertAzul3DEvent(t *testing.T) {
	tm := time.Now()

	azul3DEvent := azul3d.ButtonEvent{
		T:     tm,
		Key:   azul3d.A,
		State: azul3d.Up,
		Raw:   0,
	}

	buttonEvent := keyboard.ButtonEvent{
		T:     tm,
		Key:   keyboard.A,
		State: keyboard.Up,
	}

	convertedEvent := fromAzul3DEvent(azul3DEvent)

	if !reflect.DeepEqual(buttonEvent, convertedEvent) {
		t.Error("Failed to convert Azul3D button event")
	}
}

func TestCanCheckForAzul3DButtonEvent(t *testing.T) {
	var azul3DEvent interface{}

	azul3DEvent = azul3d.ButtonEvent{
		T:     time.Now(),
		Key:   azul3d.A,
		State: azul3d.Up,
		Raw:   0,
	}

	if !isAzul3DButtonEvent(azul3DEvent) {
		t.Error("Failed to check for Azul3D button event")
	}
}

func TestCanExtractAzul3DButtonEvent(t *testing.T) {
	var azul3DEvent interface{}

	azul3DEvent = azul3d.ButtonEvent{
		T:     time.Now(),
		Key:   azul3d.A,
		State: azul3d.Up,
		Raw:   0,
	}

	extractedEvent, ok := extractAzul3DButtonEvent(azul3DEvent)
	if !ok {
		t.Error("Failed to extract Azul3D button event")
	}

	if !reflect.DeepEqual(azul3DEvent, extractedEvent) {
		t.Error("Failed to extract Azul3D button event")
	}
}
