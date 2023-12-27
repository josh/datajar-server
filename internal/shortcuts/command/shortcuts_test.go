package command

import (
	"testing"
)

// Depends on Shortcut named "Test" that outputs 42
func TestRunShortcut(t *testing.T) {
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
	_, err := RunShortcut("DefinitelyDoesNotExist", "")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
