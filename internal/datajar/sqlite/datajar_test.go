//go:build darwin

package sqlite

import (
	"context"
	"os"
	"testing"
)

func TestFetchStore(t *testing.T) {
	ctx := context.TODO()

	_, err := os.Stat(StorePath)
	if err != nil {
		t.Skip("DataJar.sqlite does not exist")
	}

	output, err := FetchStore(ctx)
	if err != nil {
		t.Error(err)
	} else if len(output) == 0 {
		t.Error("output dictionary is empty")
	}
}
