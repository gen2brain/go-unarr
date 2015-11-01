package unarr

// #include <unarr.h>
// #cgo LDFLAGS: -lunarr
import "C"
import (
	"errors"
	"io"
	"path/filepath"
	"strings"
	"time"
	"unsafe"
)

// Archive represents unarr archive
type Archive struct {
	stream  *C.ar_stream  // C stream struct
	archive *C.ar_archive // C archive struct
}

// NewArchive returns new unarr Archive
func NewArchive(path string) (a *Archive, err error) {
	a = new(Archive)

	a.stream = C.ar_open_file(C.CString(path))
	if a.stream == nil {
		err = errors.New("unarr: File not found")
		return
	}

	var deflatedOnly bool = false
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".xps" || ext == ".epub" {
		// XPS and EPUB do not support non-Deflate compression methods by specification
		deflatedOnly = true
	}

	a.archive = C.ar_open_rar_archive(a.stream)
	if a.archive == nil {
		a.archive = C.ar_open_zip_archive(a.stream, C.bool(deflatedOnly))
	}
	if a.archive == nil {
		a.archive = C.ar_open_7z_archive(a.stream)
	}
	if a.archive == nil {
		a.archive = C.ar_open_tar_archive(a.stream)
	}

	if a.archive == nil {
		err = errors.New("unarr: No valid RAR, ZIP, 7Z or TAR archive")
	}

	return
}

// Entry reads the next archive entry
// io.EOF is returned when there is no more to be read from the archive
func (a *Archive) Entry() error {
	r := bool(C.ar_parse_entry(a.archive))
	if !r {
		e := bool(C.ar_at_eof(a.archive))
		if e {
			return io.EOF
		} else {
			return errors.New("unarr: Failed to parse entry")
		}
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
	r := bool(C.ar_parse_entry_for(a.archive, C.CString(name)))
	if !r {
		return errors.New("unarr: Entry not found")
	}

	return nil
}

// Read tries to read 'b' bytes into buffer, advancing the read offset pointer
// returns the actual number of bytes read
func (a *Archive) Read(b []byte) (n int, err error) {
	r := bool(C.ar_entry_uncompress(a.archive, unsafe.Pointer(&b[0]), C.size_t(cap(b))))

	n = len(b)
	if !r || n == 0 {
		err = io.EOF
	}

	return
}

// Size returns the total size of uncompressed data of the current entry
func (a *Archive) Size() int {
	return int(C.ar_entry_get_size(a.archive))
}

// Offset returns the stream offset of the current entry
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

// Close closes the underlying unarr archive
func (a *Archive) Close() {
	C.ar_close_archive(a.archive)
	C.ar_close(a.stream)
}
