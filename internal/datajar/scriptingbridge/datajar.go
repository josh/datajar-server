package scriptingbridge

import (
	"encoding/json"

	"github.com/josh/datajar-server/internal/shortcuts/scriptingbridge"
)

func FetchStore() (map[string]interface{}, error) {
	output, err := scriptingbridge.RunShortcut("Get Data Jar Store", "")
	if err != nil {
		return nil, err
	}
	result := output[0].(map[string]interface{})
	return result, nil
}

type shortcutInput struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func SetStoreValue(key string, value interface{}) error {
	input := shortcutInput{
		Key:   key,
		Value: value,
	}
	inputData, err := json.Marshal(input)
	if err != nil {
		return err
	}

	_, err = scriptingbridge.RunShortcut("Set Data Jar Value", string(inputData))
	return err
}
