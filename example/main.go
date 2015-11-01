package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"go-unarr"
)

func extractArchive(filename string, dest string) {
	archive, err := unarr.NewArchive(filename)
	if err != nil {
		fmt.Printf("Error NewArchive: %s\n", err)
		return
	}
	defer archive.Close()

	for {
		err := archive.Entry()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Printf("Error Next: %s\n", err)
				break
			}
		}

		name := archive.Name()
		size := archive.Size()
		mtime := archive.ModTime()

		fmt.Printf("%d\t%s %s\t%s\n", size, mtime.Format("02-01-06"), mtime.Format("15:04"), name)

		buf := make([]byte, size)
		for size > 0 {
			n, err := archive.Read(buf)
			if err != nil && err != io.EOF {
				break
			}
			size -= n
		}

		if size > 0 {
			fmt.Printf("Error Read\n")
			continue
		}

		dirname := filepath.Join(dest, filepath.Dir(name))
		os.MkdirAll(dirname, 0755)

		err = ioutil.WriteFile(filepath.Join(dirname, filepath.Base(name)), buf, 0644)
		if err != nil {
			fmt.Printf("Error WriteFile: %s\n", err)
			continue
		}
	}
}

func main() {
	tmpdir, _ := ioutil.TempDir(os.TempDir(), "unarr")
	defer os.RemoveAll(tmpdir)

	for _, filename := range os.Args[1:] {
		extractArchive(filename, tmpdir)
	}
}
