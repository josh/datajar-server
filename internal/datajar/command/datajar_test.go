package command

import (
	"testing"
)

// Depends on Shortcut named "Get Data Jar Store" that outputs 42
func TestFetchStore(t *testing.T) {
	output, err := FetchStore()
	if err != nil {
		t.Error(err)
	} else if len(output) == 0 {
		t.Error("output dictionary is empty")
	}
}

// Depends on Shortcut named "Set Data Jar Value" that accepts input
func TestSetStoreValue(t *testing.T) {
	err := SetStoreValue("foo", 42)
	if err != nil {
		t.Error(err)
	}

	err = SetStoreValue("foo", nil)
	if err != nil {
		t.Error(err)
	}
}
