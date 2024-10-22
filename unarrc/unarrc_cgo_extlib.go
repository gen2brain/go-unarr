//go:build extlib

package unarrc

/*
#cgo !pkgconfig LDFLAGS: -lunarr
#cgo pkgconfig,!static pkg-config: libunarr
#cgo pkgconfig,static pkg-config: --static libunarr
*/
import "C"
