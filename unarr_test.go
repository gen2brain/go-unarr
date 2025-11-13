package unarr

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gen2brain/go-unarr/unarrc"
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

func TestMultiVolumeRar(t *testing.T) {
	// Look for any multi-volume RAR archive in testdata/multivolume/
	// Users can place their own test data here for testing
	multivolDir := filepath.Join("testdata", "multivolume")
	
	// Check if directory exists
	if _, err := os.Stat(multivolDir); os.IsNotExist(err) {
		t.Skip("Multi-volume test directory not available (testdata/multivolume/)")
		return
	}

	// Find first .rar file in directory
	entries, err := os.ReadDir(multivolDir)
	if err != nil {
		t.Skip("Could not read multivolume directory")
		return
	}

	var rarPath string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".rar" {
			rarPath = filepath.Join(multivolDir, entry.Name())
			break
		}
	}

	if rarPath == "" {
		t.Skip("No .rar file found in testdata/multivolume/")
		return
	}

	a, err := NewArchive(rarPath)
	if err != nil {
		t.Error(err)
		return
	}
	defer a.Close()

	// Extract to temporary directory
	tmpdir, err := os.MkdirTemp(os.TempDir(), "unarr-multivolume")
	if err != nil {
		t.Error(err)
		return
	}
	defer os.RemoveAll(tmpdir)

	contents, err := a.Extract(tmpdir)
	if err != nil {
		t.Error(err)
		return
	}

	if len(contents) == 0 {
		t.Error("No files extracted from multi-volume archive")
		return
	}

	// Verify at least one file was extracted and has content
	for _, filename := range contents {
		extractedPath := filepath.Join(tmpdir, filename)
		info, err := os.Stat(extractedPath)
		if err != nil {
			t.Errorf("Failed to stat extracted file %s: %v", filename, err)
			continue
		}
		
		if info.Size() == 0 {
			t.Errorf("Extracted file %s has zero size", filename)
		}
	}

	t.Logf("Successfully extracted %d file(s) from multi-volume RAR archive", len(contents))
}

func TestVolumePathNaming(t *testing.T) {
	basePath := "/tmp/archive"

	tests := []struct {
		volumeNum int64
		expected  string
	}{
		// First volume
		{-1, "/tmp/archive.rar"},
		
		// Volumes 0-99 (.r00 to .r99)
		{0, "/tmp/archive.r00"},
		{1, "/tmp/archive.r01"},
		{50, "/tmp/archive.r50"},
		{99, "/tmp/archive.r99"},
		
		// Volumes 100-199 (.s00 to .s99)
		{100, "/tmp/archive.s00"},
		{101, "/tmp/archive.s01"},
		{150, "/tmp/archive.s50"},
		{199, "/tmp/archive.s99"},
		
		// Volumes 200-299 (.t00 to .t99)
		{200, "/tmp/archive.t00"},
		{201, "/tmp/archive.t01"},
		{250, "/tmp/archive.t50"},
		{299, "/tmp/archive.t99"},
		
		// Volumes 300-399 (.u00 to .u99)
		{300, "/tmp/archive.u00"},
		{399, "/tmp/archive.u99"},
		
		// Volumes 400-499 (.v00 to .v99)
		{400, "/tmp/archive.v00"},
		{499, "/tmp/archive.v99"},
		
		// Volumes 500-599 (.w00 to .w99)
		{500, "/tmp/archive.w00"},
		{599, "/tmp/archive.w99"},
		
		// Volumes 600-699 (.x00 to .x99)
		{600, "/tmp/archive.x00"},
		{699, "/tmp/archive.x99"},
		
		// Volumes 700-799 (.y00 to .y99)
		{700, "/tmp/archive.y00"},
		{799, "/tmp/archive.y99"},
		
		// Volumes 800-899 (.z00 to .z99)
		{800, "/tmp/archive.z00"},
		{899, "/tmp/archive.z99"},
		
		// Edge cases at boundaries
		{98, "/tmp/archive.r98"},
		{99, "/tmp/archive.r99"},
		{100, "/tmp/archive.s00"},
		{198, "/tmp/archive.s98"},
		{199, "/tmp/archive.s99"},
		{200, "/tmp/archive.t00"},
	}

	for _, tt := range tests {
		result := unarrc.TestGetVolumePath(basePath, tt.volumeNum)
		if result == "" {
			t.Errorf("TestGetVolumePath(%d) returned empty string", tt.volumeNum)
			continue
		}
		
		if result != tt.expected {
			t.Errorf("TestGetVolumePath(%d) = %s, want %s", tt.volumeNum, result, tt.expected)
		}
	}
}

func TestVolumePathMaxLimit(t *testing.T) {
	basePath := "/tmp/archive"

	// Test maximum valid volume (899 = .z99)
	result := unarrc.TestGetVolumePath(basePath, 899)
	if result == "" {
		t.Error("TestGetVolumePath(899) returned empty, but should be valid (.z99)")
	} else {
		expected := "/tmp/archive.z99"
		if result != expected {
			t.Errorf("TestGetVolumePath(899) = %s, want %s", result, expected)
		}
	}

	// Test volume beyond maximum (900 should fail)
	result = unarrc.TestGetVolumePath(basePath, 900)
	if result != "" {
		t.Errorf("TestGetVolumePath(900) should return empty (too high), but got %s", result)
	}
}

func TestVolumePathAlphabetCoverage(t *testing.T) {
	basePath := "/tmp/test"

	// Test that each letter from r to z is covered
	alphabet := "rstuvwxyz"
	for i, letter := range alphabet {
		volumeNum := int64(i * 100) // 0, 100, 200, 300, ..., 800
		result := unarrc.TestGetVolumePath(basePath, volumeNum)
		if result == "" {
			t.Errorf("TestGetVolumePath(%d) returned empty for letter '%c'", volumeNum, letter)
			continue
		}
		
		expected := "/tmp/test." + string(letter) + "00"
		if result != expected {
			t.Errorf("TestGetVolumePath(%d) = %s, want %s (letter '%c')", volumeNum, result, expected, letter)
		}
	}
}
