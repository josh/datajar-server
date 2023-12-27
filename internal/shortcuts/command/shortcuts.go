package command

import (
	"encoding/json"
	"os/exec"
)

func RunShortcut(name string) (interface{}, error) {
	cmd := exec.Command("shortcuts", "run", name)
	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var data interface{}
	err = json.Unmarshal(stdout, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
