//go:build darwin && cgo

package scriptingbridge

import (
	"context"

	"github.com/josh/datajar-server/internal/datajar/shortcuts"
	"github.com/josh/datajar-server/internal/shortcuts/scriptingbridge"
)

func FetchStore(ctx context.Context) (map[string]interface{}, error) {
	output, err := scriptingbridge.RunShortcut(ctx, "Get Data Jar Store", "")
	if err != nil {
		return nil, err
	}
	result := output[0].(map[string]interface{})
	return result, nil
}

func SetStoreValue(ctx context.Context, key string, value interface{}) error {
	input, err := shortcuts.PrepareShortcutInput(key, value)
	if err != nil {
		return err
	}
	_, err = scriptingbridge.RunShortcut(ctx, "Set Data Jar Value", input)
	return err
}
