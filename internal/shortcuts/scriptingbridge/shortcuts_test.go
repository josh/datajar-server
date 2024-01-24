//go:build darwin && cgo

package scriptingbridge

import "testing"

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
	} else if len(output) != 1 {
		t.Errorf("expected 1, got %v", len(output))
	} else if _, ok := output[0].(string); !ok {
		t.Errorf("expected string, got %T", output[0])
	} else if output[0] != "42" {
		t.Errorf("expected 42, got %v", output[0])
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
