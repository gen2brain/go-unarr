package unarrc

/*
#cgo !debug CFLAGS: -DNDEBUG
#cgo CFLAGS: -DHAVE_7Z -DHAVE_ZLIB -DHAVE_BZIP2 -D_7ZIP_PPMD_SUPPPORT
#cgo CFLAGS: -Iexternal/lzma -Iexternal/zlib -Iexternal/bzip2 -Iexternal/unarr -fomit-frame-pointer -flto
*/
import "C"
