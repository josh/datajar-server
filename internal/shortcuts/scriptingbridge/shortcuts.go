package scriptingbridge

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework ScriptingBridge
#include <stdlib.h>
#import "ShortcutsHelper.h"
*/
import "C"

import (
	"encoding/json"
	"fmt"
	"unsafe"
)

func RunShortcut(name string) ([]interface{}, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cResult := C.ShortcutResult{}
	defer C.free(unsafe.Pointer(cResult.bytes))
	defer C.free(unsafe.Pointer(cResult.err))

	C.runShortcut(cName, &cResult)

	if cResult.err != nil {
		return nil, fmt.Errorf(C.GoString(cResult.err))
	}

	jsonBytes := C.GoBytes(cResult.bytes, C.int(cResult.length))

	if len(jsonBytes) == 0 {
		return nil, nil
	}

	var data []interface{}
	err := json.Unmarshal(jsonBytes, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
