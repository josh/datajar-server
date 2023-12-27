package datajar

import (
	"reflect"
	"testing"

	"github.com/josh/datajar-server/internal/datajar/command"
	"github.com/josh/datajar-server/internal/datajar/scriptingbridge"
	"github.com/josh/datajar-server/internal/datajar/sqlite"
)

func TestFetchStore(t *testing.T) {
	commandOutput, err := command.FetchStore()
	if err != nil {
		t.Error(err)
	}

	scriptingbridgeOutput, err := scriptingbridge.FetchStore()
	if err != nil {
		t.Error(err)
	}

	sqliteOutput, err := sqlite.FetchStore()
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(commandOutput, sqliteOutput) {
		t.Errorf("command and sqlite output are different")
	}
	if !reflect.DeepEqual(commandOutput, scriptingbridgeOutput) {
		t.Errorf("command and scriptingbridge output are different")
	}
	if !reflect.DeepEqual(scriptingbridgeOutput, sqliteOutput) {
		t.Errorf("scriptingbridge and sqlite output are different")
	}
}
