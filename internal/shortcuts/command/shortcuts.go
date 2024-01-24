//go:build darwin

package command

import (
	"encoding/json"
	"os/exec"
	"strings"
	"sync"
)

var mutex = &sync.Mutex{}

func HasShortcut(name string) (bool, error) {
	mutex.Lock()
	defer mutex.Unlock()

	cmd := exec.Command("shortcuts", "list")
	stdout, err := cmd.Output()
	if err != nil {
		return false, err
	}

	names := strings.Split(string(stdout), "\n")
	for _, n := range names {
		if n == name {
			return true, nil
		}
	}

	return false, nil
}

func RunShortcut(name string, input string) (interface{}, error) {
	mutex.Lock()
	defer mutex.Unlock()

	cmd := exec.Command("shortcuts", "run", name)
	cmd.Stdin = strings.NewReader(input)

	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	if len(stdout) == 0 {
		return nil, nil
	}

	var data interface{}
	err = json.Unmarshal(stdout, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
