package main

import (
	"fmt"
	"os"

	"github.com/gen2brain/go-unarr"
)

func main() {
	if 3 != len(os.Args) {
		fmt.Fprintf(os.Stderr, "Usage: unarr <archive>    <local directory>\n")
		fmt.Fprintf(os.Stderr, "Usage: unarr ./example.7z /tmp/example/\n")
		os.Exit(1)
		return
	}

	archive := os.Args[1]
	extract := os.Args[2]

	a, err := unarr.NewArchive(archive)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
		return
	}
	defer a.Close()

	filenames, err := a.Extract(extract)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
		return
	}
	for _, name := range filenames {
		fmt.Println(name)
	}
}
