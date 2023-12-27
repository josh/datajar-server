package internal

import (
	"path/filepath"
	"strings"
)

type Capabilities struct {
	Read  []string `json:"read"`
	Write []string `json:"write"`
}

func CanAccessPath(path string, caps []Capabilities, accessType string) bool {
	for _, cap := range caps {
		var accessList []string
		if accessType == "read" {
			accessList = cap.Read
		} else if accessType == "write" {
			accessList = cap.Write
		}

		for _, access := range accessList {
			if access == "*" {
				return true
			}

			allowed, _ := filepath.Match("/"+access, path)
			if allowed {
				return true
			}
		}
	}
	return false
}

func GetPath(store map[string]interface{}, path string) interface{} {
	if path == "" || path == "/" {
		return store
	}

	parts := strings.Split(strings.Trim(path, "/"), "/")

	current := interface{}(store)
	for _, part := range parts {
		if cmap, ok := current.(map[string]interface{}); ok {
			current = cmap[part]
		} else {
			return nil
		}
	}
	return current
}
