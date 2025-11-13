//go:build !extlib

package unarrc

/*
#include "external/unarr/common/stream-multivolume-test.c"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// TestGetVolumePath is a test helper to verify volume path naming
func TestGetVolumePath(basePath string, volumeNum int64) string {
	cBasePath := C.CString(basePath)
	defer C.free(unsafe.Pointer(cBasePath))
	
	result := C.test_get_volume_path(cBasePath, C.longlong(volumeNum))
	if result == nil {
		return ""
	}
	
	goResult := C.GoString(result)
	C.test_free_path(result)
	
	return goResult
}
