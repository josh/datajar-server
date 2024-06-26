//go:build darwin && cgo

package datajar

import (
	"context"
	"reflect"
	"testing"

	"github.com/josh/datajar-server/internal/datajar/command"
	"github.com/josh/datajar-server/internal/datajar/scriptingbridge"
	"github.com/josh/datajar-server/internal/datajar/sqlite"
	shortcuts "github.com/josh/datajar-server/internal/shortcuts/command"
)

func TestFetchStore(t *testing.T) {
	ctx := context.TODO()

	if ok, err := shortcuts.HasShortcut(ctx, "Get Data Jar Store"); err != nil || !ok {
		t.Skip("shortcut not found")
	}

	commandOutput, err := command.FetchStore(ctx)
	if err != nil {
		t.Error(err)
	}

	scriptingbridgeOutput, err := scriptingbridge.FetchStore(ctx)
	if err != nil {
		t.Error(err)
	}

	sqliteOutput, err := sqlite.FetchStore(ctx)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(commandOutput, sqliteOutput) {
		t.Errorf("command and sqlite output are different:\n%v\n%v", commandOutput, sqliteOutput)
	}
	if !reflect.DeepEqual(commandOutput, scriptingbridgeOutput) {
		t.Errorf("command and scriptingbridge output are different:\n%v\n%v", commandOutput, scriptingbridgeOutput)
	}
	if !reflect.DeepEqual(scriptingbridgeOutput, sqliteOutput) {
		t.Errorf("scriptingbridge and sqlite output are different:\n%v\n%v", scriptingbridgeOutput, sqliteOutput)
	}
}
