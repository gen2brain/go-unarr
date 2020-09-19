package main

import (
	"fmt"
	"os"

	"github.com/gen2brain/go-unarr"
)

// version info as per https://goreleaser.com/customization/build/
var (
	version = "v0.0.0-LOCAL"
	commit  = "g0000000"
	date    = "0000-00-00T00:00:00Z"
)

func main() {
	if 2 == len(os.Args) && "version" == os.Args[1] {
		fmt.Printf("Foobar %s %s (%s)\n", version, commit, date)
	}

	if 3 != len(os.Args) {
		fmt.Fprintf(os.Stderr, "Usage:   unarr <archive>    <local directory>\n")
		fmt.Fprintf(os.Stderr, "Example: unarr ./example.7z /tmp/example/\n")
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
