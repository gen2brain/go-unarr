package unarr

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var (
	exts  = []string{"zip", "rar", "7z", "tar"}
	files = []string{"test.zip", "test.rar", "test.7z", "test.tar"}
)

func TestNewArchive(t *testing.T) {
	for _, f := range files {
		a, err := NewArchive(filepath.Join("testdata", f))
		if err != nil {
			t.Error(err)
		}

		a.Close()
	}
}

func TestNewArchiveFromMemory(t *testing.T) {
	for _, f := range files {
		d, err := ioutil.ReadFile(filepath.Join("testdata", f))
		if err != nil {
			t.Error(err)
		}

		a, err := NewArchiveFromMemory(d)
		if err != nil {
			t.Error(err)
		}

		a.Close()
	}
}

func TestNewArchiveFromReader(t *testing.T) {
	for _, f := range files {
		d, err := os.Open(filepath.Join("testdata", f))
		if err != nil {
			t.Error(err)
		}

		a, err := NewArchiveFromReader(d)
		if err != nil {
			t.Error(err)
		}

		a.Close()
	}
}

func TestEntry(t *testing.T) {
	for _, e := range exts {
		a, err := NewArchive(filepath.Join("testdata", "test."+e))
		if err != nil {
			t.Error(err)
		}

		err = a.Entry()
		if err != nil {
			t.Error(err)
		}

		a.Close()
	}
}

func TestEntryAt(t *testing.T) {
	for _, e := range exts {
		a, err := NewArchive(filepath.Join("testdata", "test."+e))
		if err != nil {
			t.Error(err)
		}

		err = a.Entry()
		if err != nil {
			t.Error(err)
		}

		err = a.EntryAt(a.Offset())
		if err != nil {
			t.Error(err)
		}

		a.Close()
	}
}

func TestEntryFor(t *testing.T) {
	for _, e := range exts {
		a, err := NewArchive(filepath.Join("testdata", "test."+e))
		if err != nil {
			t.Error(err)
		}

		err = a.EntryFor(e)
		if err != nil {
			t.Error(err)
		}

		data, err := a.ReadAll()
		if err != nil {
			t.Error(err)
		}

		if strings.TrimSpace(string(data)) != "unarr" {
			t.Error("Invalid data")
		}

		a.Close()
	}
}

func TestRead(t *testing.T) {
	for _, e := range exts {
		a, err := NewArchive(filepath.Join("testdata", "test."+e))
		if err != nil {
			t.Error(err)
		}

		err = a.Entry()
		if err != nil {
			t.Error(err)
		}

		data := make([]byte, 5)
		n, err := a.Read(data)
		if err != nil && err != io.EOF {
			t.Error(err)
		}

		if n != 5 {
			t.Error("Read failed")
		}

		if strings.TrimSpace(string(data)) != "unarr" {
			t.Error("Invalid data")
		}

		a.Close()
	}
}

func TestSeek(t *testing.T) {
	for _, e := range exts {
		a, err := NewArchive(filepath.Join("testdata", "test."+e))
		if err != nil {
			t.Error(err)
		}

		off, err := a.Seek(3, 0)
		if err != nil {
			t.Error(err)
		}

		if off != 3 {
			t.Error("Seek failed")
		}

		a.Close()
	}
}

func TestClose(t *testing.T) {
	for _, e := range exts {
		a, err := NewArchive(filepath.Join("testdata", "test."+e))
		if err != nil {
			t.Error(err)
		}

		err = a.Close()
		if err != nil {
			t.Error(err)
		}
	}
}

func TestSize(t *testing.T) {
	for _, f := range files {
		a, err := NewArchive(filepath.Join("testdata", f))
		if err != nil {
			t.Error(err)
		}

		err = a.Entry()
		if err != nil {
			t.Error(err)
		}

		if a.Size() != 6 {
			t.Error("Invalid size")
		}

		a.Close()
	}
}

func TestOffset(t *testing.T) {
	a, err := NewArchive(filepath.Join("testdata", "test.rar"))
	if err != nil {
		t.Error(err)
	}

	err = a.Entry()
	if err != nil {
		t.Error(err)
	}

	if a.Offset() != 20 {
		t.Error("Invalid offset")
	}

	a.Close()
}

func TestName(t *testing.T) {
	for _, e := range exts {
		a, err := NewArchive(filepath.Join("testdata", "test."+e))
		if err != nil {
			t.Error(err)
		}

		err = a.Entry()
		if err != nil {
			t.Error(err)
		}

		if a.Name() != e {
			t.Error("Invalid name")
		}

		a.Close()
	}
}

func TestModTime(t *testing.T) {
	for _, f := range files {
		a, err := NewArchive(filepath.Join("testdata", f))
		if err != nil {
			t.Error(err)
		}

		err = a.Entry()
		if err != nil {
			t.Error(err)
		}

		mtime := a.ModTime()
		year, month, day := mtime.Date()

		if year != 2016 || month != time.November || day != 22 {
			t.Error("Invalid time")
		}

		a.Close()
	}
}

func TestReadAll(t *testing.T) {
	for _, e := range exts {
		a, err := NewArchive(filepath.Join("testdata", "test."+e))
		if err != nil {
			t.Error(err)
		}

		err = a.Entry()
		if err != nil {
			t.Error(err)
		}

		data, err := a.ReadAll()
		if err != nil {
			t.Error(err)
		}

		if strings.TrimSpace(string(data)) != "unarr" {
			t.Error("Invalid data")
		}

		a.Close()
	}
}

func TestReadAllFromMemory(t *testing.T) {
	for _, e := range exts {
		m, err := ioutil.ReadFile(filepath.Join("testdata", "test."+e))
		if err != nil {
			t.Error(err)
		}

		a, err := NewArchiveFromMemory(m)
		if err != nil {
			t.Error(err)
		}

		err = a.Entry()
		if err != nil {
			t.Error(err)
		}

		data, err := a.ReadAll()
		if err != nil {
			t.Error(err)
		}

		if strings.TrimSpace(string(data)) != "unarr" {
			t.Error("Invalid data")
		}

		a.Close()
	}
}

func TestExtract(t *testing.T) {
	tmpdir, err := ioutil.TempDir(os.TempDir(), "unarr")
	if err != nil {
		t.Error(err)
	}

	defer os.RemoveAll(tmpdir)

	for _, f := range files {
		a, err := NewArchive(filepath.Join("testdata", f))
		if err != nil {
			t.Error(err)
		}

		a.Extract(tmpdir)
		a.Close()
	}
}

func TestExtractFromMemory(t *testing.T) {
	tmpdir, err := ioutil.TempDir(os.TempDir(), "unarr")
	if err != nil {
		t.Error(err)
	}

	defer os.RemoveAll(tmpdir)

	for _, f := range files {
		data, err := ioutil.ReadFile(filepath.Join("testdata", f))
		if err != nil {
			t.Error(err)
		}

		a, err := NewArchiveFromMemory(data)
		if err != nil {
			t.Error(err)
		}

		err = a.Extract(tmpdir)
		if err != nil {
			t.Error(err)
		}

		a.Close()
	}
}

func TestList(t *testing.T) {
	for _, f := range files {
		a, err := NewArchive(filepath.Join("testdata", f))
		if err != nil {
			t.Error(err)
		}

		list, err := a.List()
		if err != nil {
			t.Error(err)
		}

		if len(list) != 1 {
			t.Error("List failed")
		}

		a.Close()
	}
}

func TestCorrupted(t *testing.T) {
	a, err := NewArchive(filepath.Join("testdata", "test_corrupted.zip"))
	if err != nil {
		t.Error(err)
	}

	err = a.Entry()
	if err == nil {
		t.Error("Corrupted entry parsed")
	}

	a.Close()
}

func TestNotValid(t *testing.T) {
	_, err := NewArchive(filepath.Join("testdata", "test_notvalid.zip"))
	if err == nil {
		t.Error("Not valid archive opened")
	}
}

func TestExtractCorrupted(t *testing.T) {
	tmpdir, err := ioutil.TempDir(os.TempDir(), "unarr")
	if err != nil {
		t.Error(err)
	}

	defer os.RemoveAll(tmpdir)

	a, err := NewArchive(filepath.Join("testdata", "test_corrupted.zip"))
	if err != nil {
		t.Error(err)
	}

	err = a.Extract(tmpdir)
	if err == nil {
		t.Error("Corrupted archive extracted")
	}

	a.Close()
}

func TestListCorrupted(t *testing.T) {
	a, err := NewArchive(filepath.Join("testdata", "test_corrupted.zip"))
	if err != nil {
		t.Error(err)
	}

	_, err = a.List()
	if err == nil {
		t.Error("Corrupted archive listed")
	}

	a.Close()
}

func TestReadAllWrongSize(t *testing.T) {
	a, err := NewArchive(filepath.Join("testdata", "test_wrongsize.zip"))
	if err != nil {
		t.Error(err)
	}

	err = a.Entry()
	if err != nil {
		t.Error(err)
	}

	_, err = a.ReadAll()
	if err == nil {
		t.Error("Wrong size archive readed")
	}

	a.Close()
}
