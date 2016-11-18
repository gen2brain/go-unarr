#!/usr/bin/env bash

mkdir -p build && cd build

git clone https://github.com/zeniko/unarr && cd unarr
mkdir lzma920 && cd lzma920 && curl -L http://www.7-zip.org/a/lzma920.tar.bz2 | tar -xjvp && cd ..
curl -L http://zlib.net/zlib-1.2.8.tar.gz | tar -xzvp
curl -L http://www.bzip.org/1.0.6/bzip2-1.0.6.tar.gz | tar -xzvp
curl -L https://gist.githubusercontent.com/gen2brain/89fe506863be3fb139e8/raw/8783a7d81e22ad84944d146c5e33beab6dffc641/unarr-makefile.patch | patch -p1
CFLAGS="-DHAVE_7Z -DHAVE_ZLIB -DHAVE_BZIP2 -I./lzma920/C -I./zlib-1.2.8 -I./bzip2-1.0.6" build=release make
cp build/release/libunarr.a /usr/lib64/ && cp unarr.h /usr/include
