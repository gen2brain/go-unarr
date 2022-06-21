// Package unarr is a decompression library for RAR, TAR, ZIP and 7z archives.
package unarr

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unsafe"

	"github.com/gen2brain/go-unarr/unarrc"
)

// Archive represents unarr archive
type Archive struct {
	// C stream struct
	stream *unarrc.Stream
	// C archive struct
	archive *unarrc.Archive
}

// NewArchive returns new unarr Archive
func NewArchive(path string) (a *Archive, err error) {
	a = new(Archive)

	a.stream = unarrc.OpenFile(path)
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

	a.stream = unarrc.OpenMemory(unsafe.Pointer(&b[0]), uint(len(b)))
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

		err = errors.New("unarr: No valid RAR, ZIP, 7Z or TAR archive")
	}

	return
}

// Entry reads the next archive entry.
//
// io.EOF is returned when there is no more to be read from the archive.
func (a *Archive) Entry() error {
	r := unarrc.ParseEntry(a.archive)
	if !r {
		e := unarrc.AtEof(a.archive)
		if e {
			return io.EOF
		}

		return errors.New("unarr: Failed to parse entry")
	}

	return nil
}

// EntryAt reads the archive entry at the given offset
func (a *Archive) EntryAt(off int64) error {
	r := unarrc.ParseEntryAt(a.archive, off)
	if !r {
		return errors.New("unarr: Failed to parse entry at")
	}

	return nil
}

// EntryFor reads the (first) archive entry associated with the given name
func (a *Archive) EntryFor(name string) error {
	r := unarrc.ParseEntryFor(a.archive, name)
	if !r {
		return errors.New("unarr: Entry not found")
	}

	return nil
}

// Read tries to read 'b' bytes into buffer, advancing the read offset pointer.
//
// Returns the actual number of bytes read.
func (a *Archive) Read(b []byte) (n int, err error) {
	r := unarrc.EntryUncompress(a.archive, unsafe.Pointer(&b[0]), uint(len(b)))

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
	r := unarrc.Seek(a.stream, offset, whence)
	if !r {
		return 0, errors.New("unarr: Seek failed")
	}

	return int64(unarrc.Tell(a.stream)), nil
}

// Close closes the underlying unarr archive
func (a *Archive) Close() (err error) {
	unarrc.CloseArchive(a.archive)
	unarrc.Close(a.stream)

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
	var err error
	var n int

	size := a.Size()
	b := make([]byte, size)

	for size > 0 {
		n, err = a.Read(b)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
		}

		size -= n
	}

	if size > 0 {
		return nil, fmt.Errorf("unarr read failure: %w", err)
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

func toValidName(name string) string {
	p := filepath.Clean(name)
	if strings.HasPrefix(p, "/") {
		p = p[len("/"):]
	}
	for strings.HasPrefix(p, "../") {
		p = p[len("../"):]
	}
	return p
}
