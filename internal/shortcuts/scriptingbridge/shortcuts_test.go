package scriptingbridge

import "testing"

// Depends on Shortcut named "Test" that outputs 42
func TestRunShortcut(t *testing.T) {
	output, err := RunShortcut("Test")
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
	_, err := RunShortcut("DefinitelyDoesNotExist")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
