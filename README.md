go-unarr
========

Golang bindings for the [unarr](https://github.com/sumatrapdfreader/sumatrapdf/tree/master/ext/unarr) library from sumatrapdf.

unarr is a decompression library for RAR, TAR, ZIP and 7z archives.

## Install

    $ go get -d github.com/gen2brain/go-unarr
    $ cd $GOPATH/src/github.com/gen2brain/go-unarr
    $ sudo ./make.bash # compile and install libunarr static library

    $ go get github.com/gen2brain/go-unarr

## Example

    a, err := unarr.NewArchive("test.7z")
    if err != nil {
        fmt.Printf("Error NewArchive: %s\n", err)
    }
    defer a.Close()

    a.Extract(os.TempDir())

## Docs

[docs on Go Walker](https://gowalker.org/github.com/gen2brain/go-unarr)
