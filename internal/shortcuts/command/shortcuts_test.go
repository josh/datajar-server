//go:build darwin

package command

import (
	"testing"
)

// Depends on Shortcut named "Test" that outputs 42
func TestRunShortcut(t *testing.T) {
	if ok, err := HasShortcut("Test"); err != nil {
		t.Skip("skipping test; error checking for shortcut:", err)
	} else if !ok {
		t.Skip("skipping test; shortcut not found")
	}

	output, err := RunShortcut("Test", "")
	if err != nil {
		t.Errorf("error running shortcut: %s", err)
	} else if _, ok := output.(float64); !ok {
		t.Errorf("expected float64, got %T", output)
	} else if output != 42.0 {
		t.Errorf("expected 42, got %v", output)
	}
}

func TestMissingShortcut(t *testing.T) {
	if ok, err := HasShortcut("Test"); err != nil {
		t.Skip("error checking for shortcut:", err)
	} else if !ok {
		t.Skip("shortcut not found")
	}

	_, err := RunShortcut("DefinitelyDoesNotExist", "")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
