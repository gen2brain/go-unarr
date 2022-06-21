//go:build !windows
// +build !windows

package unarrc

/*
#include "unarr.h"
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"
)

// OpenFile .
func OpenFile(path string) *Stream {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	ret := C.ar_open_file(cpath)
	v := *(**Stream)(unsafe.Pointer(&ret))
	return v
}
