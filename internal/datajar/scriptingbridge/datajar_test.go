//go:build darwin && cgo

package scriptingbridge

import (
	"context"
	"testing"

	shortcuts "github.com/josh/datajar-server/internal/shortcuts/scriptingbridge"
)

// Depends on Shortcut named "Get Data Jar Store" that outputs 42
func TestFetchStore(t *testing.T) {
	ctx := context.TODO()

	if ok, err := shortcuts.HasShortcut(ctx, "Get Data Jar Store"); err != nil || !ok {
		t.Skip("shortcut not found")
	}

	output, err := FetchStore(ctx)
	if err != nil {
		t.Error(err)
	} else if len(output) == 0 {
		t.Error("output dictionary is empty")
	}
}

// Depends on Shortcut named "Set Data Jar Value" that accepts input
func TestSetStoreValue(t *testing.T) {
	ctx := context.TODO()

	if ok, err := shortcuts.HasShortcut(ctx, "Set Data Jar Value"); err != nil || !ok {
		t.Skip("shortcut not found")
	}

	err := SetStoreValue(ctx, "foo", 42)
	if err != nil {
		t.Error(err)
	}

	err = SetStoreValue(ctx, "foo", nil)
	if err != nil {
		t.Error(err)
	}
}
