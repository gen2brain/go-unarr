//go:build windows

package unarrc

/*
#include "unarr.h"
*/
import "C"

import (
	"unicode/utf16"
	"unsafe"
)

// OpenFile .
func OpenFile(path string) *Stream {
	a := utf16.Encode([]rune(path + "\x00"))
	cpath := (*C.wchar_t)(&a[0])

	ret := C.ar_open_file_w(cpath)
	v := *(**Stream)(unsafe.Pointer(&ret))
	return v
}
