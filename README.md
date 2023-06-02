# go-unarr

[![Build Status](https://github.com/gen2brain/go-unarr/actions/workflows/test.yml/badge.svg)](https://github.com/gen2brain/go-unarr/actions)
[![GoDoc](https://godoc.org/github.com/gen2brain/go-unarr?status.svg)](https://godoc.org/github.com/gen2brain/go-unarr)
[![Go Report Card](https://goreportcard.com/badge/github.com/gen2brain/go-unarr?branch=master)](https://goreportcard.com/report/github.com/gen2brain/go-unarr)

> Golang bindings for the [unarr](https://github.com/selmf/unarr) library from sumatrapdf.

`unarr` is a decompression library and CLI for RAR, TAR, ZIP and 7z archives.

## GoDoc

See <https://pkg.go.dev/github.com/gen2brain/go-unarr>

## Install CLI

```bash
go install github.com/gen2brain/go-unarr/cmd/unarr@latest
```

#### Example

```bash
unarr ./example.7z ./example/
```

#### Build

For one-off builds:

```bash
go build -o ./unarr ./cmd/unarr/*.go
```

For multi-platform cross-compile builds:

```bash
goreleaser --snapshot --skip-publish --rm-dist
```

## Library Examples

#### Install Library

```bash
go get -v github.com/gen2brain/go-unarr
```

#### Open archive

```go
a, err := unarr.NewArchive("test.7z")
if err != nil {
    panic(err)
}
defer a.Close()
```

#### Read first entry from archive

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

#### List contents of archive

```go
list, err := a.List()
if err != nil {
    panic(err)
}
```

#### Read known filename from archive

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

#### Read first 8 bytes of the entry

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

#### Read all entries from archive

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

#### Extract contents of archive to destination path

```go
_, err := a.Extract("/tmp/path")
if err != nil {
    panic(err)
}
```
