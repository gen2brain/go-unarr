package unarrc

/*
#cgo !debug CFLAGS: -DNDEBUG
#cgo CFLAGS: -std=c99 -DHAVE_7Z -DHAVE_ZLIB -DHAVE_BZIP2 -D_7ZIP_PPMD_SUPPPORT
#cgo CFLAGS: -Iexternal/lzma -Iexternal/zlib -Iexternal/bzip2 -Iexternal/unarr -fomit-frame-pointer
#cgo amd64 arm64 arm64be loong64 mips64 mips64le ppc64 ppc64le riscv64 s390x sparc64 CFLAGS: -D_FILE_OFFSET_BITS=64
*/
import "C"
