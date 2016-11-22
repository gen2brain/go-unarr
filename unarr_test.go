package unarr

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewArchive(t *testing.T) {
	a, err := NewArchive(filepath.Join("testdata", "test.zip"))
	if err != nil {
		t.Error(err)
	}
	a.Close()
}

func TestNewArchiveFromMemory(t *testing.T) {
	d, err := ioutil.ReadFile(filepath.Join("testdata", "test.zip"))
	if err != nil {
		t.Error(err)
	}

	a, err := NewArchiveFromMemory(d)
	if err != nil {
		t.Error(err)
	}

	a.Close()
}

func TestExtract(t *testing.T) {
	tmpdir, err := ioutil.TempDir(os.TempDir(), "unarr")
	if err != nil {
		t.Error(err)
	}

	defer os.RemoveAll(tmpdir)

	files := []string{"test.zip", "test.rar", "test.7z"}

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

	files := []string{"test.zip", "test.rar", "test.7z"}

	for _, f := range files {
		data, err := ioutil.ReadFile(filepath.Join("testdata", f))
		if err != nil {
			t.Error(err)
		}

		a, err := NewArchiveFromMemory(data)
		if err != nil {
			t.Error(err)
		}

		a.Extract(tmpdir)
		a.Close()
	}
}

func TestReadAll(t *testing.T) {
	exts := []string{"zip", "rar", "7z"}

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
	exts := []string{"zip", "rar", "7z"}

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

func TestEntryFor(t *testing.T) {
	exts := []string{"zip", "rar", "7z"}

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

func TestList(t *testing.T) {
	files := []string{"test.zip", "test.rar", "test.7z"}

	for _, f := range files {
		a, err := NewArchive(filepath.Join("testdata", f))
		if err != nil {
			t.Error(err)
		}

		a.List()
		a.Close()
	}
}
