package command

import (
	"github.com/josh/datajar-server/internal/shortcuts/command"
)

func FetchStore() (map[string]interface{}, error) {
	output, err := command.RunShortcut("Get Data Jar Store")
	if err != nil {
		return nil, err
	}
	result := output.(map[string]interface{})
	return result, nil
}
