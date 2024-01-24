//go:build darwin

package sqlite

import (
	"os"
	"testing"
)

func TestFetchStore(t *testing.T) {
	_, err := os.Stat(StorePath)
	if err != nil {
		t.Skip("DataJar.sqlite does not exist")
	}

	output, err := FetchStore()
	if err != nil {
		t.Error(err)
	} else if len(output) == 0 {
		t.Error("output dictionary is empty")
	}
}
