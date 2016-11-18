go-unarr
========

Golang bindings for the [unarr](https://github.com/sumatrapdfreader/sumatrapdf/tree/master/ext/unarr) library from sumatrapdf.

unarr is a decompression library for RAR, TAR, ZIP and 7z archives.

## Install

    $ go get -d github.com/gen2brain/go-unarr
    $ cd $GOPATH/src/github.com/gen2brain/go-unarr
    $ ./make.bash # compile libunarr static library

    $ go get github.com/gen2brain/go-unarr

## Example

    a, err := unarr.NewArchive("test.rar")
    if err != nil {
        fmt.Printf("Error NewArchive: %s\n", err)
    }
    defer a.Close()

    a.Extract(os.TempDir())
