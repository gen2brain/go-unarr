## go-unarr [![Build Status](https://travis-ci.org/gen2brain/go-unarr.svg?branch=master)](https://travis-ci.org/gen2brain/go-unarr) [![GoDoc](https://godoc.org/github.com/gen2brain/go-unarr?status.svg)](https://godoc.org/github.com/gen2brain/go-unarr) [![Go Report Card](https://goreportcard.com/badge/github.com/gen2brain/go-unarr)](https://goreportcard.com/report/github.com/gen2brain/go-unarr?branch=master)

Golang bindings for the [unarr](https://github.com/sumatrapdfreader/sumatrapdf/tree/master/ext/unarr) library from sumatrapdf.

unarr is a decompression library for RAR, TAR, ZIP and 7z archives.

### Installation

    go get -v github.com/gen2brain/go-unarr

### Example

##### Open archive
```go
a, err := unarr.NewArchive("test.7z")
if err != nil {
    panic(err)
}
defer a.Close()
```

##### Read first entry from archive
```go
err := a.Entry()
if err != nil {
    panic(err)
}

data, err := a.ReadAll()
if err != nil {
    panic(err)
}
```

##### Read known filename from archive
```go
err := a.EntryFor("filename.txt")
if err != nil {
    panic(err)
}

data, err := a.ReadAll()
if err != nil {
    panic(err)
}
```

##### Read all entries from archive
```go
for {
    err := a.Entry()
    if err != nil {
        if err == io.EOF {
            break
        } else {
            panic(err)
        }
    }

    data, err := a.ReadAll()
    if err != nil {
        panic(err)
    }
}
```

##### Extract contents of archive to destination path
```go
err := a.Extract("/tmp/path")
if err != nil {
    panic(err)
}
```
