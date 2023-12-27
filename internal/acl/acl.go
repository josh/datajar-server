package acl

import (
	"path/filepath"
	"strings"
)

type Capabilities struct {
	Read []string `json:"read"`
}

func CanReadPath(path string, caps []Capabilities) bool {
	for _, cap := range caps {
		for _, read := range cap.Read {
			if read == "*" {
				return true
			}

			allowed, _ := filepath.Match("/"+read, path)
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
