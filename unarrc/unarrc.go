package unarrc

/*
#include "unarr.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// Stream .
type Stream C.ar_stream

// Archive .
type Archive C.ar_archive

// OpenFile .
func OpenFile(path string) *Stream {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	__ret := C.ar_open_file(cpath)
	__v := *(**Stream)(unsafe.Pointer(&__ret))
	return __v
}

// OpenMemory .
func OpenMemory(data unsafe.Pointer, datalen uint) *Stream {
	cdata := data
	cdatalen := (C.size_t)(datalen)
	__ret := C.ar_open_memory(cdata, cdatalen)
	__v := *(**Stream)(unsafe.Pointer(&__ret))
	return __v
}

// Close .
func Close(stream *Stream) {
	cstream := (*C.ar_stream)(unsafe.Pointer(stream))
	C.ar_close(cstream)
}

// Read .
func Read(stream *Stream, buffer unsafe.Pointer, count uint) uint {
	cstream := (*C.ar_stream)(unsafe.Pointer(stream))
	cbuffer := buffer
	ccount := (C.size_t)(count)
	__ret := C.ar_read(cstream, cbuffer, ccount)
	__v := (uint)(__ret)
	return __v
}

// Seek .
func Seek(stream *Stream, offset int64, origin int) bool {
	cstream := (*C.ar_stream)(unsafe.Pointer(stream))
	coffset := (C.off64_t)(offset)
	corigin := (C.int)(origin)
	__ret := C.ar_seek(cstream, coffset, corigin)
	__v := (bool)(__ret)
	return __v
}

// Skip .
func Skip(stream *Stream, count int64) bool {
	cstream := (*C.ar_stream)(unsafe.Pointer(stream))
	ccount := (C.off64_t)(count)
	__ret := C.ar_skip(cstream, ccount)
	__v := (bool)(__ret)
	return __v
}

// Tell .
func Tell(stream *Stream) int {
	cstream := (*C.ar_stream)(unsafe.Pointer(stream))
	__ret := C.ar_tell(cstream)
	__v := (int)(__ret)
	return __v
}

// CloseArchive .
func CloseArchive(ar *Archive) {
	car := (*C.ar_archive)(unsafe.Pointer(ar))
	C.ar_close_archive(car)
}

// ParseEntry .
func ParseEntry(ar *Archive) bool {
	car := (*C.ar_archive)(unsafe.Pointer(ar))
	__ret := C.ar_parse_entry(car)
	__v := (bool)(__ret)
	return __v
}

// ParseEntryAt .
func ParseEntryAt(ar *Archive, offset int64) bool {
	car := (*C.ar_archive)(unsafe.Pointer(ar))
	coffset := (C.off64_t)(offset)
	__ret := C.ar_parse_entry_at(car, coffset)
	__v := (bool)(__ret)
	return __v
}

// ParseEntryFor .
func ParseEntryFor(ar *Archive, entryName string) bool {
	car := (*C.ar_archive)(unsafe.Pointer(ar))
	centryName := C.CString(entryName)
	defer C.free(unsafe.Pointer(centryName))
	__ret := C.ar_parse_entry_for(car, centryName)
	__v := (bool)(__ret)
	return __v
}

// AtEof .
func AtEof(ar *Archive) bool {
	car := (*C.ar_archive)(unsafe.Pointer(ar))
	__ret := C.ar_at_eof(car)
	__v := (bool)(__ret)
	return __v
}

// EntryGetName .
func EntryGetName(ar *Archive) string {
	car := (*C.ar_archive)(unsafe.Pointer(ar))
	__ret := C.ar_entry_get_name(car)
	__v := C.GoString(__ret)
	return __v
}

// EntryGetRawName .
func EntryGetRawName(ar *Archive) string {
	car := (*C.ar_archive)(unsafe.Pointer(ar))
	__ret := C.ar_entry_get_raw_name(car)
	__v := C.GoString(__ret)
	return __v
}

// EntryGetOffset .
func EntryGetOffset(ar *Archive) int {
	car := (*C.ar_archive)(unsafe.Pointer(ar))
	__ret := C.ar_entry_get_offset(car)
	__v := (int)(__ret)
	return __v
}

// EntryGetSize .
func EntryGetSize(ar *Archive) uint {
	car := (*C.ar_archive)(unsafe.Pointer(ar))
	__ret := C.ar_entry_get_size(car)
	__v := (uint)(__ret)
	return __v
}

// EntryGetFiletime .
func EntryGetFiletime(ar *Archive) int {
	car := (*C.ar_archive)(unsafe.Pointer(ar))
	__ret := C.ar_entry_get_filetime(car)
	__v := (int)(__ret)
	return __v
}

// EntryUncompress .
func EntryUncompress(ar *Archive, buffer unsafe.Pointer, count uint) bool {
	car := (*C.ar_archive)(unsafe.Pointer(ar))
	cbuffer := buffer
	ccount := (C.size_t)(count)
	__ret := C.ar_entry_uncompress(car, cbuffer, ccount)
	__v := (bool)(__ret)
	return __v
}

// GetGlobalComment .
func GetGlobalComment(ar *Archive, buffer unsafe.Pointer, count uint) uint {
	car := (*C.ar_archive)(unsafe.Pointer(ar))
	cbuffer := buffer
	ccount := (C.size_t)(count)
	__ret := C.ar_get_global_comment(car, cbuffer, ccount)
	__v := (uint)(__ret)
	return __v
}

// OpenRarArchive .
func OpenRarArchive(stream *Stream) *Archive {
	cstream := (*C.ar_stream)(unsafe.Pointer(stream))
	__ret := C.ar_open_rar_archive(cstream)
	__v := *(**Archive)(unsafe.Pointer(&__ret))
	return __v
}

// OpenTarArchive .
func OpenTarArchive(stream *Stream) *Archive {
	cstream := (*C.ar_stream)(unsafe.Pointer(stream))
	__ret := C.ar_open_tar_archive(cstream)
	__v := *(**Archive)(unsafe.Pointer(&__ret))
	return __v
}

// OpenZipArchive .
func OpenZipArchive(stream *Stream, deflatedonly bool) *Archive {
	cstream := (*C.ar_stream)(unsafe.Pointer(stream))
	cdeflatedonly := (C._Bool)(deflatedonly)
	__ret := C.ar_open_zip_archive(cstream, cdeflatedonly)
	__v := *(**Archive)(unsafe.Pointer(&__ret))
	return __v
}

// Open7zArchive .
func Open7zArchive(stream *Stream) *Archive {
	cstream := (*C.ar_stream)(unsafe.Pointer(stream))
	__ret := C.ar_open_7z_archive(cstream)
	__v := *(**Archive)(unsafe.Pointer(&__ret))
	return __v
}
