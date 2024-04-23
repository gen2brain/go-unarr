package unarrc

/*
#cgo !debug CFLAGS: -DNDEBUG -DFUZZING_BUILD_MODE_UNSAFE_FOR_PRODUCTION
#cgo CFLAGS: -std=c99 -DHAVE_7Z -DHAVE_ZLIB -DHAVE_BZIP2 -D_7ZIP_PPMD_SUPPPORT -DBZ_NO_STDIO
#cgo CFLAGS: -Iexternal/lzma -Iexternal/zlib -Iexternal/bzip2 -Iexternal/unarr -fomit-frame-pointer -Wno-typedef-redefinition -Wno-unknown-warning
#cgo amd64 arm64 arm64be loong64 mips64 mips64le ppc64 ppc64le riscv64 s390x sparc64 CFLAGS: -D_FILE_OFFSET_BITS=64

void bz_internal_error (int errcode) {};
*/
import "C"
