|Build|Status|
|---|------|
|Linux|[![Build Status](https://dev.azure.com/licorn0647/licorn/_apis/build/status/selmf.unarr?branchName=master&jobName=Linux)](https://dev.azure.com/licorn0647/licorn/_build/latest?definitionId=2&branchName=master)|
|MacOS|[![Build Status](https://dev.azure.com/licorn0647/licorn/_apis/build/status/selmf.unarr?branchName=master&jobName=MacOS)](https://dev.azure.com/licorn0647/licorn/_build/latest?definitionId=2&branchName=master)|
|Windows|[![Build Status](https://dev.azure.com/licorn0647/licorn/_apis/build/status/selmf.unarr?branchName=master&jobName=Windows)](https://dev.azure.com/licorn0647/licorn/_build/latest?definitionId=2&branchName=master)|

# (lib)unarr

**(lib)unarr** is a decompression library for RAR, TAR, ZIP and 7z* archives.

It was forked from **unarr**, which originated as a port of the RAR extraction
features from The Unarchiver project required for extracting images from comic
book archives. [Zeniko](https://github.com/zeniko/) wrote unarr as an
alternative to libarchive which didn't have support for parsing filters or
solid compression at the time.

While (lib)unarr was started with the intent of providing unarr with a proper
cmake based build system suitable for packaging and cross-platform development,
it's focus has now been extended to provide code maintenance and to continue the
development of unarr, which no longer is maintained.

## Getting started

### Prebuilt packages
[![Packaging status](https://repology.org/badge/vertical-allrepos/unarr.svg)](https://repology.org/metapackage/unarr)

#### From OBS
[.deb package](https://software.opensuse.org//download.html?project=home%3Aselmf&package=libunarr)
[.rpm package](https://software.opensuse.org//download.html?project=home%3Aselmf%3Ayacreader-rpm&package=libunarr)

### Building from source

#### Dependencies

(lib)unarr can take advantage of the following libraries if they are present:

* bzip2
* xz / libLZMA
* zlib

More information on what library is used for which purpose can be found in the
description for embedded builds.

#### CMake

```bash
mkdir build
cd build
cmake ..
make
```

... as a static library

```bash
cmake .. -DBUILD_SHARED_LIBS=OFF
```

... with 7z support (see note below)

```bash
cmake .. -DENABLE_7Z=ON
```

By default, (lib)unarr will try to detect and use system libraries like bzip2,
xz/LibLZMA and zlib. If this is undesirable, you can override this behavior by
specifying:

```bash
cmake .. -DUSE_SYSTEM_BZ2=OFF -DUSE_SYSTEM_LZMA=OFF -DUSE_SYSTEM_ZLIB=OFF
```

Install

```bash
make install
```

#### Embedded build

Make sure your compiler is C99 compatible, grab the source code, copy it into
your project and adjust your build system accordingly.

You can define the following symbols to take advantage of third party libraries:

| Symbol            | Required header | Required for (format/method)|
|-------------------|:---------------:|:----------------------------|
|HAVE_ZLIB          |     zlib.h      |  faster CRC-32 and Deflate  |
|HAVE_BZIP2         |     bzlib.h     |    ZIP / Bzip2              |
|HAVE_LIBLZMA       |     lzma.h      |    ZIP / LZMA, XZ(LZMA2)    |
|HAVE_7Z            |     7z.h        |    7Z / LZMA, LZMA2, BCJ    |
|_7ZIP_PPMD_SUPPPORT|                 |    7Z / PPMd                |

Make sure the required headers are present in the include path.

## Usage

### Examples

Check [unarr.h](unarr.h) and [unarr-test](test/main.c) to get a general feel
for the api and usage.

To build the unarr-test sample application, use:

```bash
cmake .. -DBUILD_SAMPLES=ON
```

## Limitations

Unarr was written for comic book archives, so it currently doesn't support:

* password protected archives
* self extracting archives
* split archives

### 7z support

7z support is currently limited and has to be explicitly enabled at build time.

This is due to a known performance problem in the ANSI-C based 7z extraction
code provided by the LZMA SDK that limits its usefulness for large files with
solid compression (see https://github.com/zeniko/unarr/issues/4).

Fixing this problem will require modification or replacement of the LZMA SDK
code used.
