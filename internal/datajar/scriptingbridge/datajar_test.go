package scriptingbridge

import (
	"testing"

	shortcuts "github.com/josh/datajar-server/internal/shortcuts/scriptingbridge"
)

// Depends on Shortcut named "Get Data Jar Store" that outputs 42
func TestFetchStore(t *testing.T) {
	if ok, err := shortcuts.HasShortcut("Get Data Jar Store"); err != nil || !ok {
		t.Skip("shortcut not found")
	}

	output, err := FetchStore()
	if err != nil {
		t.Error(err)
	} else if len(output) == 0 {
		t.Error("output dictionary is empty")
	}
}

// Depends on Shortcut named "Set Data Jar Value" that accepts input
func TestSetStoreValue(t *testing.T) {
	if ok, err := shortcuts.HasShortcut("Set Data Jar Value"); err != nil || !ok {
		t.Skip("shortcut not found")
	}

	err := SetStoreValue("foo", 42)
	if err != nil {
		t.Error(err)
	}

	err = SetStoreValue("foo", nil)
	if err != nil {
		t.Error(err)
	}
}
