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
	"sync"
	"unsafe"
)

var mutex = &sync.Mutex{}

func HasShortcut(name string) (bool, error) {
	mutex.Lock()
	defer mutex.Unlock()

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	ok := C.hasShortcut(cName)
	return ok == 1, nil
}

func RunShortcut(name string, input string) ([]interface{}, error) {
	mutex.Lock()
	defer mutex.Unlock()

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cInput := C.CString(input)
	defer C.free(unsafe.Pointer(cInput))

	cResult := C.ShortcutResult{}
	defer C.free(unsafe.Pointer(cResult.bytes))
	defer C.free(unsafe.Pointer(cResult.err))

	C.runShortcut(cName, cInput, &cResult)

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
