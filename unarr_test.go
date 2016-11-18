package unarr

import (
	"io/ioutil"
	"os"
	"path/filepath"
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
