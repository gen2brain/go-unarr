// Package unarr is a decompression library for RAR, TAR, ZIP and 7z archives.
package unarr

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"unsafe"

	"github.com/gen2brain/go-unarr/unarrc"
)

var (
	ErrOpenFile    = errors.New("unarr: open file failed")
	ErrOpenMemory  = errors.New("unarr: open memory failed")
	ErrOpenArchive = errors.New("unarr: no valid RAR, ZIP, 7Z or TAR archive")
	ErrEntry       = errors.New("unarr: failed to parse entry")
	ErrEntryAt     = errors.New("unarr: failed to parse entry at")
	ErrEntryFor    = errors.New("unarr: failed to parse entry for")
	ErrSeek        = errors.New("unarr: seek failed")
	ErrRead        = errors.New("unarr: read failure")
	ErrUnsafePath  = errors.New("unarr: entry path is not local")
)

const readAllPrealloc = 1 << 20

// Archive represents unarr archive
type Archive struct {
	// C stream struct
	stream *unarrc.Stream
	// C archive struct
	archive *unarrc.Archive
	// archive data, pinned while open
	data   []byte
	pinner runtime.Pinner
	// unread bytes of the current entry
	remaining int
}

// NewArchive returns new unarr Archive
func NewArchive(path string) (a *Archive, err error) {
	a = new(Archive)

	a.stream = unarrc.OpenFile(path)
	if a.stream == nil {
		err = ErrOpenFile
		return
	}

	err = a.open()

	return
}

// NewArchiveFromMemory returns new unarr Archive from byte slice.
//
// The byte slice must not be modified until the archive is closed.
func NewArchiveFromMemory(b []byte) (a *Archive, err error) {
	a = new(Archive)

	if len(b) == 0 {
		err = ErrOpenMemory
		return
	}

	a.data = b
	a.pinner.Pin(&b[0])

	a.stream = unarrc.OpenMemory(unsafe.Pointer(&b[0]), uint(len(b)))
	if a.stream == nil {
		a.Close()
		err = ErrOpenMemory
		return
	}

	err = a.open()
	if err != nil {
		a.Close()
	}

	return
}

// NewArchiveFromReader returns new unarr Archive from io.Reader
func NewArchiveFromReader(r io.Reader) (a *Archive, err error) {
	b, e := io.ReadAll(r)
	if e != nil {
		err = e
		return
	}

	a, err = NewArchiveFromMemory(b)

	return
}

// open opens archive
func (a *Archive) open() (err error) {
	a.archive = unarrc.OpenRarArchive(a.stream)
	if a.archive == nil {
		a.archive = unarrc.OpenZipArchive(a.stream, false)
	}
	if a.archive == nil {
		a.archive = unarrc.Open7zArchive(a.stream)
	}
	if a.archive == nil {
		a.archive = unarrc.OpenTarArchive(a.stream)
	}

	if a.archive == nil {
		unarrc.Close(a.stream)
		a.stream = nil

		err = ErrOpenArchive

		return
	}

	runtime.SetFinalizer(a, (*Archive).Close)

	return
}

// Entry reads the next archive entry.
//
// io.EOF is returned when there is no more to be read from the archive.
func (a *Archive) Entry() error {
	a.remaining = 0

	r := unarrc.ParseEntry(a.archive)
	if !r {
		e := unarrc.AtEof(a.archive)
		if e {
			return io.EOF
		}

		return ErrEntry
	}

	a.remaining = a.Size()

	return nil
}

// EntryAt reads the archive entry at the given offset
func (a *Archive) EntryAt(off int64) error {
	a.remaining = 0

	r := unarrc.ParseEntryAt(a.archive, off)
	if !r {
		return ErrEntryAt
	}

	a.remaining = a.Size()

	return nil
}

// EntryFor reads the (first) archive entry associated with the given name
func (a *Archive) EntryFor(name string) error {
	a.remaining = 0

	r := unarrc.ParseEntryFor(a.archive, name)
	if !r {
		return ErrEntryFor
	}

	a.remaining = a.Size()

	return nil
}

// Read reads up to len(b) bytes of the current entry into b, advancing the read offset pointer.
//
// Returns the actual number of bytes read, io.EOF at the end of the entry.
func (a *Archive) Read(b []byte) (n int, err error) {
	if a.remaining <= 0 {
		return 0, io.EOF
	}

	if len(b) == 0 {
		return 0, nil
	}

	n = min(len(b), a.remaining)
	if !unarrc.EntryUncompress(a.archive, unsafe.Pointer(&b[0]), uint(n)) {
		a.remaining = 0
		return 0, ErrRead
	}

	a.remaining -= n

	return n, nil
}

// Seek moves the read offset pointer interpreted according to whence.
//
// Returns the new offset.
func (a *Archive) Seek(offset int64, whence int) (int64, error) {
	r := unarrc.Seek(a.stream, offset, whence)
	if !r {
		return 0, ErrSeek
	}

	return int64(unarrc.Tell(a.stream)), nil
}

// Close closes the underlying unarr archive
func (a *Archive) Close() (err error) {
	if a.archive != nil {
		unarrc.CloseArchive(a.archive)
		a.archive = nil
	}

	if a.stream != nil {
		unarrc.Close(a.stream)
		a.stream = nil
	}

	a.pinner.Unpin()
	a.data = nil
	a.remaining = 0

	runtime.SetFinalizer(a, nil)

	return
}

// Size returns the total size of uncompressed data of the current entry
func (a *Archive) Size() int {
	return int(unarrc.EntryGetSize(a.archive))
}

// Offset returns the stream offset of the current entry, for use with EntryAt
func (a *Archive) Offset() int64 {
	return int64(unarrc.EntryGetOffset(a.archive))
}

// Name returns the name of the current entry as UTF-8 string
func (a *Archive) Name() string {
	return toValidName(unarrc.EntryGetName(a.archive))
}

// RawName returns the name of the current entry as raw string
func (a *Archive) RawName() string {
	return unarrc.EntryGetRawName(a.archive)
}

// ModTime returns the stored modification time of the current entry
func (a *Archive) ModTime() time.Time {
	filetime := int64(unarrc.EntryGetFiletime(a.archive))
	return time.Unix((filetime/10000000)-11644473600, 0)
}

// ReadAll reads current entry and returns data
func (a *Archive) ReadAll() ([]byte, error) {
	var buf bytes.Buffer
	buf.Grow(min(a.remaining, readAllPrealloc))

	_, err := buf.ReadFrom(a)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
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
		if !filepath.IsLocal(name) {
			err = ErrUnsafePath
			return
		}

		contents = append(contents, name)
		data, e := a.ReadAll()
		if e != nil {
			err = e
			return
		}

		dirname := filepath.Join(path, filepath.Dir(name))
		e = os.MkdirAll(dirname, 0755)
		if e != nil {
			err = e
			return
		}

		e = os.WriteFile(filepath.Join(dirname, filepath.Base(name)), data, 0644)
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

func toValidName(name string) string {
	p := path.Clean("/" + strings.ReplaceAll(name, `\`, "/"))

	return filepath.FromSlash(strings.TrimPrefix(p, "/"))
}
