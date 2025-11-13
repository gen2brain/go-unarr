# go-unarr

[![Build Status](https://github.com/gen2brain/go-unarr/actions/workflows/test.yml/badge.svg)](https://github.com/gen2brain/go-unarr/actions)
[![GoDoc](https://godoc.org/github.com/gen2brain/go-unarr?status.svg)](https://godoc.org/github.com/gen2brain/go-unarr)
[![Go Report Card](https://goreportcard.com/badge/github.com/gen2brain/go-unarr?branch=master)](https://goreportcard.com/report/github.com/gen2brain/go-unarr)

> Golang bindings for the [unarr](https://github.com/selmf/unarr) library from sumatrapdf.

`unarr` is a decompression library and CLI for RAR, TAR, ZIP and 7z archives.

## Features

- Extract RAR, TAR, ZIP, and 7z archives
- **Multi-volume RAR archive support** (.rar, .r00-.r99, .s00-.s99, .t00-.t99, .u00-.u99, .v00-.v99, .w00-.w99, .x00-.x99, .y00-.y99, .z00-.z99)
- List archive contents
- Read individual files from archives
- Pure Go API with CGO bindings to C libraries

## GoDoc

See <https://pkg.go.dev/github.com/gen2brain/go-unarr>

## Build tags

* `extlib` - use external libunarr library
* `pkgconfig` - enable pkg-config (used with `extlib`)
* `static` - use static library (used with `pkgconfig`)

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

#### Extract multi-volume RAR archive

```go
// Provide the path to the first volume (.rar file)
// The library will automatically read from subsequent volumes:
//   .r00, .r01, ..., .r99 (volumes 0-99)
//   .s00, .s01, ..., .s99 (volumes 100-199)
//   .t00, .t01, ..., .t99 (volumes 200-299)
//   ... continuing through alphabet ...
//   .z00, .z01, ..., .z99 (volumes 800-899, up to 900 total volumes)
a, err := unarr.NewArchive("archive.rar")
if err != nil {
    panic(err)
}
defer a.Close()

_, err = a.Extract("/tmp/path")
if err != nil {
    panic(err)
}
```

## Testing

Run all tests:

```bash
go test -v
```

### Multi-Volume Archive Testing (Optional)

To test multi-volume RAR support, place your multi-volume test files in `testdata/multivolume/`:

```bash
mkdir -p testdata/multivolume
# Copy your multi-volume files to testdata/multivolume/
# Supports: .rar, .r00-.r99, .s00-.s99, .t00-.t99, etc.
go test -v -run TestMultiVolumeRar
```

The test will automatically detect and test any `.rar` file found in that directory.
