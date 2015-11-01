go-unarr
========

Golang bindings for the [unarr](https://github.com/sumatrapdfreader/sumatrapdf/tree/master/ext/unarr) library from sumatrapdf.

unarr is a decompression library for RAR, TAR, ZIP and 7z archives.

Install
-------

    $ git clone https://github.com/zeniko/unarr && cd unarr
    $ mkdir lzma920 && cd lzma920 && curl -L http://www.7-zip.org/a/lzma920.tar.bz2 | tar -xjvp && cd ..
    $ curl -L https://gist.githubusercontent.com/gen2brain/89fe506863be3fb139e8/raw/826eb3d03d64eda2119b16863bcb5c0ccdc7c473/unarr-makefile.patch | patch -p1
    $ CFLAGS="-DHAVE_7Z -I./lzma920/C" make
    $ cp build/debug/libunarr.a /usr/lib64/ && cp unarr.h /usr/include

    $ go get github.com/gen2brain/go-unarr
