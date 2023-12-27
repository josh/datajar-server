package scriptingbridge

import (
	"github.com/josh/datajar-server/internal/shortcuts/scriptingbridge"
)

func FetchStore() (map[string]interface{}, error) {
	output, err := scriptingbridge.RunShortcut("Get Data Jar Store")
	if err != nil {
		return nil, err
	}
	result := output[0].(map[string]interface{})
	return result, nil
}
