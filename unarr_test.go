package unarr

import (
	"archive/zip"
	"bytes"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
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
		d, err := os.ReadFile(filepath.Join("testdata", f))
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
		m, err := os.ReadFile(filepath.Join("testdata", "test."+e))
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
	tmpdir, err := os.MkdirTemp(os.TempDir(), "unarr")
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
	tmpdir, err := os.MkdirTemp(os.TempDir(), "unarr")
	if err != nil {
		t.Error(err)
	}

	defer os.RemoveAll(tmpdir)

	for _, f := range files {
		data, err := os.ReadFile(filepath.Join("testdata", f))
		if err != nil {
			t.Error(err)
		}

		a, err := NewArchiveFromMemory(data)
		if err != nil {
			t.Error(err)
		}

		_, err = a.Extract(tmpdir)
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
	tmpdir, err := os.MkdirTemp(os.TempDir(), "unarr")
	if err != nil {
		t.Error(err)
	}

	defer os.RemoveAll(tmpdir)

	a, err := NewArchive(filepath.Join("testdata", "test_corrupted.zip"))
	if err != nil {
		t.Error(err)
	}

	_, err = a.Extract(tmpdir)
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

func TestSubdir(t *testing.T) {
	a, err := NewArchive(filepath.Join("testdata", "test_subdir.zip"))
	if err != nil {
		t.Error(err)
	}

	err = a.EntryFor("test/test.txt")
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

func TestEmpty(t *testing.T) {
	for _, b := range [][]byte{nil, {}} {
		_, err := NewArchiveFromMemory(b)
		if err != ErrOpenMemory {
			t.Errorf("expected ErrOpenMemory, got %v", err)
		}
	}

	_, err := NewArchiveFromReader(bytes.NewReader(nil))
	if err != ErrOpenMemory {
		t.Errorf("expected ErrOpenMemory, got %v", err)
	}
}

func openDropped(t *testing.T, f string) *Archive {
	d, err := os.ReadFile(filepath.Join("testdata", f))
	if err != nil {
		t.Fatal(err)
	}

	a, err := NewArchiveFromMemory(d)
	if err != nil {
		t.Fatal(err)
	}

	return a
}

func TestMemoryOwnership(t *testing.T) {
	for _, f := range files {
		a := openDropped(t, f)

		for i := 0; i < 5; i++ {
			runtime.GC()
		}

		err := a.Entry()
		if err != nil {
			t.Errorf("%s: %v", f, err)
		}

		data, err := a.ReadAll()
		if err != nil {
			t.Errorf("%s: %v", f, err)
		}

		if strings.TrimSpace(string(data)) != "unarr" {
			t.Errorf("%s: invalid data %q", f, data)
		}

		a.Close()
	}
}

func TestNoClose(t *testing.T) {
	for i := 0; i < 10; i++ {
		for _, f := range files {
			openDropped(t, f)
		}

		runtime.GC()
		time.Sleep(10 * time.Millisecond)
	}
}

func TestReadBuffer(t *testing.T) {
	for _, f := range files {
		a, err := NewArchive(filepath.Join("testdata", f))
		if err != nil {
			t.Fatal(err)
		}

		err = a.Entry()
		if err != nil {
			t.Fatal(err)
		}

		n, err := a.Read(nil)
		if n != 0 || err != nil {
			t.Errorf("%s: empty buffer: n=%d err=%v", f, n, err)
		}

		size := a.Size()
		data := make([]byte, size+100)

		n, err = a.Read(data)
		if n != size || err != nil {
			t.Errorf("%s: n=%d err=%v, expected n=%d", f, n, err, size)
		}

		if strings.TrimSpace(string(data[:n])) != "unarr" {
			t.Errorf("%s: invalid data %q", f, data[:n])
		}

		n, err = a.Read(data)
		if n != 0 || err != io.EOF {
			t.Errorf("%s: after end: n=%d err=%v", f, n, err)
		}

		a.Close()
	}
}

func TestCloseTwice(t *testing.T) {
	d, err := os.ReadFile(filepath.Join("testdata", "test.zip"))
	if err != nil {
		t.Fatal(err)
	}

	a, err := NewArchiveFromMemory(d)
	if err != nil {
		t.Fatal(err)
	}

	a.Close()
	a.Close()
}

func TestExtractUnsafePath(t *testing.T) {
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	for _, name := range []string{"../evil1.txt", "/evil2.txt", "a/../../evil3.txt", `..\evil4.txt`} {
		f, err := w.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		f.Write([]byte("evil"))
	}
	w.Close()

	tmpdir := t.TempDir()
	dest := filepath.Join(tmpdir, "dest")

	a, err := NewArchiveFromMemory(buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	defer a.Close()

	_, err = a.Extract(dest)
	if err != nil && err != ErrUnsafePath {
		t.Fatal(err)
	}

	filepath.WalkDir(tmpdir, func(p string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() && !strings.HasPrefix(p, dest+string(filepath.Separator)) {
			t.Errorf("file outside destination: %s", p)
		}

		return nil
	})
}

func TestValidName(t *testing.T) {
	for name, want := range map[string]string{
		"../a.txt":          "a.txt",
		"/a.txt":            "a.txt",
		"a/../../b.txt":     "b.txt",
		`..\a.txt`:          "a.txt",
		`dir\..\..\a.txt`:   "a.txt",
		"dir/sub/a.txt":     filepath.Join("dir", "sub", "a.txt"),
		"./dir//a.txt":      filepath.Join("dir", "a.txt"),
		"../../../../a.txt": "a.txt",
	} {
		if got := toValidName(name); got != want {
			t.Errorf("toValidName(%q) = %q, want %q", name, got, want)
		}
	}
}
