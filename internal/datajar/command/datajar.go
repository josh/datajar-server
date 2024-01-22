package command

import (
	"github.com/josh/datajar-server/internal/datajar/shortcuts"
	"github.com/josh/datajar-server/internal/shortcuts/command"
)

func FetchStore() (map[string]interface{}, error) {
	output, err := command.RunShortcut("Get Data Jar Store", "")
	if err != nil {
		return nil, err
	}
	result := output.(map[string]interface{})
	return result, nil
}

func SetStoreValue(key string, value interface{}) error {
	input, err := shortcuts.PrepareShortcutInput(key, value)
	if err != nil {
		return err
	}
	_, err = command.RunShortcut("Set Data Jar Value", input)
	return err
}
