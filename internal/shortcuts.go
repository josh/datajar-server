package internal

import (
	"encoding/json"
	"strings"
)

type ShortcutInput struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func ConvertToJSONPath(requestPath string) string {
	jsonPath := requestPath
	jsonPath = strings.Join(strings.Split(strings.Trim(requestPath, "/"), "/"), ".")
	return jsonPath
}

func PrepareShortcutInput(key string, value interface{}) (string, error) {
	input := ShortcutInput{
		Key:   ConvertToJSONPath(key),
		Value: value,
	}
	inputData, err := json.Marshal(input)
	if err != nil {
		return "", err
	}
	return string(inputData), nil
}
