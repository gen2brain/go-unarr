#!/bin/bash

# This cross-compiled for all supported cross-compile targets with goreleaser
# from MacOS, like so:
#     goreleaser --snapshot --skip-publish --rm-dist

# It can also be cross-compiled manually, like so:

mkdir -p dist

export CGO_ENABLED=1
printf "Build for Darwin (macOS) ...            "
GOARCH=amd64 GOOS=darwin go build -o dist/unarr_darwin_amd64 -ldflags '-w' cmd/unarr/unarr.go
echo "dist/unarr_darwin_amd64"

# Install x86_64-w64-mingw32-gcc:
#     brew install mingw-w64
printf "Build for Windows ...                   "
GOARCH=amd64 GOOS=windows CC=x86_64-w64-mingw32-gcc go build -o dist/unarr.exe -ldflags '-w' cmd/unarr/unarr.go
echo "dist/unarr.exe"

# Install x86_64-linux-musl-gcc:
#     brew install FiloSottile/musl-cross/musl-cross
# can also use '-w -linkmode external -extldflags "-static"'
printf "Build for Linux ...                     "
GOARCH=amd64 GOOS=linux   CC=x86_64-linux-musl-gcc go build -o dist/unarr_linux_musl_amd64 --ldflags '-w -extldflags "-static"' cmd/unarr/unarr.go
echo "dist/unarr_linux_musl_amd64"

# Install arm* toolchains:
#     brew install FiloSottile/musl-cross/musl-cross --with-arm --with-arm-hf --with-aarch64
printf "Build for Raspberry Pi (64-bit ARM) ... "
GOARCH=arm64 GOOS=linux   CC=aarch64-linux-musl-gcc go build -o dist/unarr_linux_musl_aarch64 --ldflags '-w -extldflags "-static"' cmd/unarr/unarr.go
echo "dist/unarr_linux_musl_aarch64"
