package server

import (
	"errors"
	"strings"
)

// Given a HTTP request path, traverse the store and return the value at that path.
func GetValueByPath(store map[string]interface{}, requestPath string) (interface{}, error) {
	parts := removeEmptyStrings(strings.Split(strings.Trim(requestPath, "/"), "/"))

	if len(parts) == 0 {
		return store, nil
	}

	current := interface{}(store)
	for _, part := range parts {
		cmap, ok := current.(map[string]interface{})
		if !ok {
			return nil, errors.New("path not found")
		}
		value, exists := cmap[part]
		if !exists {
			return nil, errors.New("path not found")
		}
		current = value
	}

	return current, nil
}

func removeEmptyStrings(input []string) []string {
	result := make([]string, 0, len(input))
	for _, str := range input {
		if str != "" {
			result = append(result, str)
		}
	}
	return result
}
