/* Copyright 2015 the unarr project authors (see AUTHORS file).
   License: LGPLv3 */

/* adapted from https://code.google.com/p/theunarchiver/source/browse/XADMaster/XADRARParser.m */

#include "rar.h"
#include <stdint.h>

static inline uint8_t uint8le(unsigned char *data) { return data[0]; }
static inline uint16_t uint16le(unsigned char *data) { return data[0] | data[1] << 8; }
static inline uint32_t uint32le(unsigned char *data) { return data[0] | data[1] << 8 | data[2] << 16 | (uint32_t)data[3] << 24; }

bool rar_parse_header(ar_archive *ar, struct rar_header *header)
{
    unsigned char header_data[7];
    size_t read = ar_read(ar->stream, header_data, sizeof(header_data));
    if (read == 0) {
        ar->at_eof = true;
        return false;
    }
    if (read < sizeof(header_data))
        return false;

    header->crc = uint16le(header_data + 0);
    header->type = uint8le(header_data + 2);
    header->flags = uint16le(header_data + 3);
    header->size = uint16le(header_data + 5);

    header->datasize = 0;
    if ((header->flags & LHD_LONG_BLOCK) || header->type == 0x74) {
        unsigned char size_data[4];
        if (!(header->flags & LHD_LONG_BLOCK))
            log("File header without LHD_LONG_BLOCK set");
        read += ar_read(ar->stream, size_data, sizeof(size_data));
        if (read < sizeof(header_data) + sizeof(size_data))
            return false;
        header->datasize = uint32le(size_data);
    }

    if (header->size < read) {
        warn("Invalid header size %d", header->size);
        return false;
    }

    return true;
}

bool rar_check_header_crc(ar_archive *ar)
{
    unsigned char buffer[256];
    uint16_t crc16, size;
    uint32_t crc32;

    if (!ar_seek(ar->stream, ar->entry_offset, SEEK_SET))
        return false;
    if (ar_read(ar->stream, buffer, 7) != 7)
        return false;

    crc16 = uint16le(buffer + 0);
    size = uint16le(buffer + 5);
    if (size < 7)
        return false;
    size -= 7;

    crc32 = ar_crc32(0, buffer + 2, 5);
    while (size > 0) {
        if (ar_read(ar->stream, buffer, smin(size, sizeof(buffer))) != smin(size, sizeof(buffer)))
            return false;
        crc32 = ar_crc32(crc32, buffer, smin(size, sizeof(buffer)));
        size -= (uint16_t)smin(size, sizeof(buffer));
    }
    return (crc32 & 0xFFFF) == crc16;
}

bool rar_parse_header_entry(ar_archive_rar *rar, struct rar_header *header, struct rar_entry *entry)
{
    unsigned char data[21];
    if (ar_read(rar->super.stream, data, sizeof(data)) != sizeof(data))
        return false;

    entry->size = uint32le(data + 0);
    entry->os = uint8le(data + 4);
    entry->crc = uint32le(data + 5);
    entry->dosdate = uint32le(data + 9);
    entry->version = uint8le(data + 13);
    entry->method = uint8le(data + 14);
    entry->namelen = uint16le(data + 15);
    entry->attrs = uint32le(data + 17);
    if ((header->flags & LHD_LARGE)) {
        unsigned char more_data[8];
        if (ar_read(rar->super.stream, more_data, sizeof(more_data)) != sizeof(more_data))
            return false;
        header->datasize += (uint64_t)uint32le(more_data + 0);
        entry->size += (uint64_t)uint32le(more_data + 4);
    }
    if (!ar_skip(rar->super.stream, entry->namelen))
        return false;
    if ((header->flags & LHD_SALT)) {
        log("Skipping LHD_SALT");
        ar_skip(rar->super.stream, 8);
    }

    rar->entry.version = entry->version;
    rar->entry.method = entry->method;
    rar->entry.crc = entry->crc;
    rar->entry.header_size = header->size;
    rar->entry.solid = entry->version < 20 ? (rar->archive_flags & MHD_SOLID) : (header->flags & LHD_SOLID);
    free(rar->entry.name);
    rar->entry.name = NULL;

    return true;
}

/* Decode RAR encoded Unicode filename to UTF-8 */
static char *rar_conv_unicode_to_utf8(const char *data, uint16_t len)
{
    // Encoded UTF-16 data starts after the \0 separator
    const uint8_t *in = (const uint8_t *)data + strlen(data) + 1;
    const uint8_t *end_in = (const uint8_t *)data + len;

    // Allocate 3 bytes per character
    char *str = calloc(len + 1, 3);
    if (!str) {
        return NULL;
    }
    char *out = str;
    char *end_out = str + (len * 3);

    // No encoded Unicode data
    if ((end_in - in) <= 1) {
        memcpy(str, data, len);
        return str;
    }

    uint8_t highbyte = *in++;
    uint8_t flagbyte = 0;
    uint8_t flagbits = 0;
    uint8_t name_index = 0;

    while ((in < end_in) && (out < end_out)) {
        if (flagbits == 0) {
            flagbyte = *in++;
            flagbits = 8;
        }
        flagbits -= 2;
        switch ((flagbyte >> flagbits) & 3) {
            case 0:
            // Use the next byte as Unicode code point
            if (in + 1 > end_in) {
                goto fail;
            }
            out += ar_conv_rune_to_utf8(*in++, out, end_out - out);
            name_index++;
            break;
        case 1:
            // Combine stored high byte with next input byte
            if (in + 1 > end_in) {
                goto fail;
            }
            out += ar_conv_rune_to_utf8(((uint16_t)highbyte << 8) | *in++, out, end_out - out);
            name_index++;
            break;
        case 2:
            // Read next two bytes as little-endian 16-bit value
            if (in + 2 > end_in) {
                goto fail;
            }
            out += ar_conv_rune_to_utf8(((uint16_t)*(in + 1) << 8) | *in, out, end_out - out);
            in += 2;
            name_index++;
            break;
        case 3:
            {
                // Reconstruct a series of Unicode code points from filename with optional correction value
                if (in + 1 > end_in) {
                    goto fail;
                }
                uint8_t length = *in++;
                if ((length & 0x80)) {
                    if (in + 1 > end_in) {
                        goto fail;
                    }
                    uint8_t correction = *in++;
                    for (uint8_t i = 0; i < (length & 0x7F) + 2; i++) {
                        if (name_index >= len) {
                            goto fail;
                        }
                        out += ar_conv_rune_to_utf8(((uint16_t)highbyte << 8) |
                                                    (data[name_index] + (correction & 0xFF)), out, end_out - out);
                        name_index++;
                    }
                }
                else {
                    for (uint8_t i = 0; i < (length & 0x7F) + 2; i++) {
                        if (name_index >= len) {
                            goto fail;
                        }
                        out += ar_conv_rune_to_utf8(data[name_index], out, end_out - out);
                        name_index++;
                    }
                }
                break;
            }
        }
    }
    return str;

fail:
    free(str);
    warn("Failed to parse RAR entry name");
    return NULL;
}

const char *rar_get_name(ar_archive *ar, bool raw)
{
    if (raw)
        return NULL;

    ar_archive_rar *rar = (ar_archive_rar *)ar;
    if (!rar->entry.name) {
        unsigned char data[21];
        uint16_t namelen;
        char *name;

        struct rar_header header;
        if (!ar_seek(ar->stream, ar->entry_offset, SEEK_SET))
            return NULL;
        if (!rar_parse_header(ar, &header))
            return NULL;
        if (ar_read(ar->stream, data, sizeof(data)) != sizeof(data))
            return NULL;
        if ((header.flags & LHD_LARGE) && !ar_skip(ar->stream, 8))
            return NULL;

        namelen = uint16le(data + 15);
        name = malloc(namelen + 1);
        if (!name || ar_read(ar->stream, name, namelen) != namelen) {
            free(name);
            return NULL;
        }
        name[namelen] = '\0';

        if (!(header.flags & LHD_UNICODE)) {
            rar->entry.name = ar_conv_dos_to_utf8(name);
            free(name);
        }
        /* no null separator - filename is UTF-8 */
        else if (namelen == strlen(name)) {
            rar->entry.name = name;
        }
        /* filename + null separator + encoded unicode data */
        else {
            rar->entry.name = rar_conv_unicode_to_utf8(name, namelen);
            free(name);
        }
        /* normalize path separators */
        if (rar->entry.name) {
            char *p = rar->entry.name;
            while ((p = strchr(p, '\\')) != NULL) {
                *p = '/';
            }
        }

        if (!ar_seek(ar->stream, ar->entry_offset + rar->entry.header_size, SEEK_SET))
            warn("Couldn't seek back to the end of the entry header");
    }
    return rar->entry.name;
}
