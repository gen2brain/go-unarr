## go-unarr
[![TravisCI Build Status](https://travis-ci.org/gen2brain/go-unarr.svg?branch=master)](https://travis-ci.org/gen2brain/go-unarr) 
[![AppVeyor Build Status](https://ci.appveyor.com/api/projects/status/vawnh8s1j3w6la9t?svg=true)](https://ci.appveyor.com/project/gen2brain/go-unarr)
[![GoDoc](https://godoc.org/github.com/gen2brain/go-unarr?status.svg)](https://godoc.org/github.com/gen2brain/go-unarr) 
[![Go Report Card](https://goreportcard.com/badge/github.com/gen2brain/go-unarr?branch=master)](https://goreportcard.com/report/github.com/gen2brain/go-unarr) 
<!--[![Go Cover](http://gocover.io/_badge/github.com/gen2brain/go-unarr)](http://gocover.io/github.com/gen2brain/go-unarr)-->

Golang bindings for the [unarr](https://github.com/sumatrapdfreader/sumatrapdf/tree/master/ext/unarr) library from sumatrapdf.

unarr is a decompression library for RAR, TAR, ZIP and 7z archives.

### Installation

    go get -v github.com/gen2brain/go-unarr

### Examples

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
##### List contents of archive
```go
list, err := a.List()
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

##### Read first 8 bytes of the entry
```go
err := a.Entry()
if err != nil {
    panic(err)
}

data := make([]byte, 8)

n, err := a.Read(data)
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
