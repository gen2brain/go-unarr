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
	ret := C.ar_open_file(cpath)
	v := *(**Stream)(unsafe.Pointer(&ret))
	return v
}

// OpenMemory .
func OpenMemory(data unsafe.Pointer, datalen uint) *Stream {
	cdata := data
	cdatalen := (C.size_t)(datalen)
	ret := C.ar_open_memory(cdata, cdatalen)
	v := *(**Stream)(unsafe.Pointer(&ret))
	return v
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
	ret := C.ar_read(cstream, cbuffer, ccount)
	v := (uint)(ret)
	return v
}

// Seek .
func Seek(stream *Stream, offset int64, origin int) bool {
	cstream := (*C.ar_stream)(unsafe.Pointer(stream))
	coffset := (C.off64_t)(offset)
	corigin := (C.int)(origin)
	ret := C.ar_seek(cstream, coffset, corigin)
	v := (bool)(ret)
	return v
}

// Skip .
func Skip(stream *Stream, count int64) bool {
	cstream := (*C.ar_stream)(unsafe.Pointer(stream))
	ccount := (C.off64_t)(count)
	ret := C.ar_skip(cstream, ccount)
	v := (bool)(ret)
	return v
}

// Tell .
func Tell(stream *Stream) int {
	cstream := (*C.ar_stream)(unsafe.Pointer(stream))
	ret := C.ar_tell(cstream)
	v := (int)(ret)
	return v
}

// CloseArchive .
func CloseArchive(ar *Archive) {
	car := (*C.ar_archive)(unsafe.Pointer(ar))
	C.ar_close_archive(car)
}

// ParseEntry .
func ParseEntry(ar *Archive) bool {
	car := (*C.ar_archive)(unsafe.Pointer(ar))
	ret := C.ar_parse_entry(car)
	v := (bool)(ret)
	return v
}

// ParseEntryAt .
func ParseEntryAt(ar *Archive, offset int64) bool {
	car := (*C.ar_archive)(unsafe.Pointer(ar))
	coffset := (C.off64_t)(offset)
	ret := C.ar_parse_entry_at(car, coffset)
	v := (bool)(ret)
	return v
}

// ParseEntryFor .
func ParseEntryFor(ar *Archive, entryName string) bool {
	car := (*C.ar_archive)(unsafe.Pointer(ar))
	centryName := C.CString(entryName)
	defer C.free(unsafe.Pointer(centryName))
	ret := C.ar_parse_entry_for(car, centryName)
	v := (bool)(ret)
	return v
}

// AtEof .
func AtEof(ar *Archive) bool {
	car := (*C.ar_archive)(unsafe.Pointer(ar))
	ret := C.ar_at_eof(car)
	v := (bool)(ret)
	return v
}

// EntryGetName .
func EntryGetName(ar *Archive) string {
	car := (*C.ar_archive)(unsafe.Pointer(ar))
	ret := C.ar_entry_get_name(car)
	v := C.GoString(ret)
	return v
}

// EntryGetRawName .
func EntryGetRawName(ar *Archive) string {
	car := (*C.ar_archive)(unsafe.Pointer(ar))
	ret := C.ar_entry_get_raw_name(car)
	v := C.GoString(ret)
	return v
}

// EntryGetOffset .
func EntryGetOffset(ar *Archive) int {
	car := (*C.ar_archive)(unsafe.Pointer(ar))
	ret := C.ar_entry_get_offset(car)
	v := (int)(ret)
	return v
}

// EntryGetSize .
func EntryGetSize(ar *Archive) uint {
	car := (*C.ar_archive)(unsafe.Pointer(ar))
	ret := C.ar_entry_get_size(car)
	v := (uint)(ret)
	return v
}

// EntryGetFiletime .
func EntryGetFiletime(ar *Archive) int {
	car := (*C.ar_archive)(unsafe.Pointer(ar))
	ret := C.ar_entry_get_filetime(car)
	v := (int)(ret)
	return v
}

// EntryUncompress .
func EntryUncompress(ar *Archive, buffer unsafe.Pointer, count uint) bool {
	car := (*C.ar_archive)(unsafe.Pointer(ar))
	cbuffer := buffer
	ccount := (C.size_t)(count)
	ret := C.ar_entry_uncompress(car, cbuffer, ccount)
	v := (bool)(ret)
	return v
}

// GetGlobalComment .
func GetGlobalComment(ar *Archive, buffer unsafe.Pointer, count uint) uint {
	car := (*C.ar_archive)(unsafe.Pointer(ar))
	cbuffer := buffer
	ccount := (C.size_t)(count)
	ret := C.ar_get_global_comment(car, cbuffer, ccount)
	v := (uint)(ret)
	return v
}

// OpenRarArchive .
func OpenRarArchive(stream *Stream) *Archive {
	cstream := (*C.ar_stream)(unsafe.Pointer(stream))
	ret := C.ar_open_rar_archive(cstream)
	v := *(**Archive)(unsafe.Pointer(&ret))
	return v
}

// OpenTarArchive .
func OpenTarArchive(stream *Stream) *Archive {
	cstream := (*C.ar_stream)(unsafe.Pointer(stream))
	ret := C.ar_open_tar_archive(cstream)
	v := *(**Archive)(unsafe.Pointer(&ret))
	return v
}

// OpenZipArchive .
func OpenZipArchive(stream *Stream, deflatedonly bool) *Archive {
	cstream := (*C.ar_stream)(unsafe.Pointer(stream))
	cdeflatedonly := (C._Bool)(deflatedonly)
	ret := C.ar_open_zip_archive(cstream, cdeflatedonly)
	v := *(**Archive)(unsafe.Pointer(&ret))
	return v
}

// Open7zArchive .
func Open7zArchive(stream *Stream) *Archive {
	cstream := (*C.ar_stream)(unsafe.Pointer(stream))
	ret := C.ar_open_7z_archive(cstream)
	v := *(**Archive)(unsafe.Pointer(&ret))
	return v
}
