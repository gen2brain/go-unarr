// Package unarr is a decompression library for RAR, TAR, ZIP and 7z archives.
package unarr

/*
#include <stdlib.h>
#include <unarr.h>
*/
import "C"

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
	"unsafe"
)

// Archive represents unarr archive
type Archive struct {
	// C stream struct
	stream *C.ar_stream
	// C archive struct
	archive *C.ar_archive
}

// NewArchive returns new unarr Archive
func NewArchive(path string) (a *Archive, err error) {
	a = new(Archive)

	p := C.CString(path)
	defer C.free(unsafe.Pointer(p))

	a.stream = C.ar_open_file(p)
	if a.stream == nil {
		err = errors.New("unarr: File not found")
		return
	}

	err = a.open()

	return
}

// NewArchiveFromMemory returns new unarr Archive from byte slice
func NewArchiveFromMemory(b []byte) (a *Archive, err error) {
	a = new(Archive)

	a.stream = C.ar_open_memory(unsafe.Pointer(&b[0]), C.size_t(cap(b)))
	if a.stream == nil {
		err = errors.New("unarr: Open memory failed")
		return
	}

	err = a.open()

	return
}

// NewArchiveFromReader returns new unarr Archive from io.Reader
func NewArchiveFromReader(r io.Reader) (a *Archive, err error) {
	b, e := ioutil.ReadAll(r)
	if e != nil {
		err = e
		return
	}

	a, err = NewArchiveFromMemory(b)

	return
}

// open opens archive
func (a *Archive) open() (err error) {
	a.archive = C.ar_open_rar_archive(a.stream)
	if a.archive == nil {
		a.archive = C.ar_open_zip_archive(a.stream, C.bool(false))
	}
	if a.archive == nil {
		a.archive = C.ar_open_7z_archive(a.stream)
	}
	if a.archive == nil {
		a.archive = C.ar_open_tar_archive(a.stream)
	}

	if a.archive == nil {
		C.ar_close(a.stream)

		err = errors.New("unarr: No valid RAR, ZIP, 7Z or TAR archive")
	}

	return
}

// Entry reads the next archive entry.
//
// io.EOF is returned when there is no more to be read from the archive.
func (a *Archive) Entry() error {
	r := bool(C.ar_parse_entry(a.archive))
	if !r {
		e := bool(C.ar_at_eof(a.archive))
		if e {
			return io.EOF
		}

		return errors.New("unarr: Failed to parse entry")
	}

	return nil
}

// EntryAt reads the archive entry at the given offset
func (a *Archive) EntryAt(off int64) error {
	r := bool(C.ar_parse_entry_at(a.archive, C.off64_t(off)))
	if !r {
		return errors.New("unarr: Failed to parse entry at")
	}

	return nil
}

// EntryFor reads the (first) archive entry associated with the given name
func (a *Archive) EntryFor(name string) error {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))

	r := bool(C.ar_parse_entry_for(a.archive, n))
	if !r {
		return errors.New("unarr: Entry not found")
	}

	return nil
}

// Read tries to read 'b' bytes into buffer, advancing the read offset pointer.
//
// Returns the actual number of bytes read.
func (a *Archive) Read(b []byte) (n int, err error) {
	r := bool(C.ar_entry_uncompress(a.archive, unsafe.Pointer(&b[0]), C.size_t(cap(b))))

	n = len(b)
	if !r || n == 0 {
		err = io.EOF
	}

	return
}

// Seek moves the read offset pointer interpreted according to whence.
//
// Returns the new offset.
func (a *Archive) Seek(offset int64, whence int) (int64, error) {
	r := bool(C.ar_seek(a.stream, C.off64_t(offset), C.int(whence)))
	if !r {
		return 0, errors.New("unarr: Seek failed")
	}

	return int64(C.ar_tell(a.stream)), nil
}

// Close closes the underlying unarr archive
func (a *Archive) Close() (err error) {
	C.ar_close_archive(a.archive)
	C.ar_close(a.stream)

	return
}

// Size returns the total size of uncompressed data of the current entry
func (a *Archive) Size() int {
	return int(C.ar_entry_get_size(a.archive))
}

// Offset returns the stream offset of the current entry, for use with EntryAt
func (a *Archive) Offset() int64 {
	return int64(C.ar_entry_get_offset(a.archive))
}

// Name returns the name of the current entry
func (a *Archive) Name() string {
	return C.GoString(C.ar_entry_get_name(a.archive))
}

// ModTime returns the stored modification time of the current entry
func (a *Archive) ModTime() time.Time {
	filetime := int64(C.ar_entry_get_filetime(a.archive))
	return time.Unix((filetime/10000000)-11644473600, 0)
}

// ReadAll reads current entry and returns data
func (a *Archive) ReadAll() ([]byte, error) {
	size := a.Size()
	read := size

	b := make([]byte, size)

	for size > 0 {
		n, err := a.Read(b)
		if err != nil && err != io.EOF {
			return nil, err
		}

		size -= n

		if err != io.EOF {
			read -= n
		}
	}

	if read > 0 {
		err := errors.New("unarr: Error Read")
		return nil, err
	}

	return b, nil
}

// Extract extracts archive to destination path
func (a *Archive) Extract(path string) (contents []string, err error) {
	for {
		e := a.Entry()
		if e != nil {
			if e == io.EOF {
				break
			}

			err = e
			return
		}

		name := a.Name()
		contents = append(contents, name)
		data, e := a.ReadAll()
		if e != nil {
			err = e
			return
		}

		dirname := filepath.Join(path, filepath.Dir(name))
		os.MkdirAll(dirname, 0755)

		e = ioutil.WriteFile(filepath.Join(dirname, filepath.Base(name)), data, 0644)
		if e != nil {
			err = e
			return
		}
	}

	return
}

// List lists the contents of archive
func (a *Archive) List() (contents []string, err error) {
	for {
		e := a.Entry()
		if e != nil {
			if e == io.EOF {
				break
			}

			err = e
			return
		}

		name := a.Name()
		contents = append(contents, name)
	}

	return
}
