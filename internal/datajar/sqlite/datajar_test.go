package sqlite

import (
	"testing"
)

func TestFetchStore(t *testing.T) {
	output, err := FetchStore()
	if err != nil {
		t.Error(err)
	} else if len(output) == 0 {
		t.Error("output dictionary is empty")
	}
}
