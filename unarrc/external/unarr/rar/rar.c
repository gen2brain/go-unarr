/* Copyright 2015 the unarr project authors (see AUTHORS file).
   License: LGPLv3 */

#include "rar.h"

static void rar_close(ar_archive *ar)
{
    ar_archive_rar *rar = (ar_archive_rar *)ar;
    free(rar->entry.name);
    rar_clear_uncompress(&rar->uncomp);
}

static bool rar_parse_entry(ar_archive *ar, off64_t offset)
{
    ar_archive_rar *rar = (ar_archive_rar *)ar;
    struct rar_header header;
    struct rar_entry entry;
    bool out_of_order = offset != ar->entry_offset_next;

    if (!ar_seek(ar->stream, offset, SEEK_SET)) {
        warn("Couldn't seek to offset %" PRIi64, offset);
        return false;
    }

    for (;;) {
        ar->entry_offset = ar_tell(ar->stream);
        ar->entry_size_uncompressed = 0;

        if (!rar_parse_header(ar, &header))
            return false;

        ar->entry_offset_next = ar->entry_offset + header.size + header.datasize;
        if (ar->entry_offset_next < ar->entry_offset + header.size) {
            warn("Integer overflow due to overly large data size");
            return false;
        }

        switch (header.type) {
        case TYPE_MAIN_HEADER:
            if ((header.flags & MHD_PASSWORD)) {
                warn("Encrypted archives aren't supported");
                return false;
            }
            ar_skip(ar->stream, 6 /* reserved data */);
            if ((header.flags & MHD_ENCRYPTVER)) {
                log("MHD_ENCRYPTVER is set");
                ar_skip(ar->stream, 1);
            }
            if ((header.flags & MHD_COMMENT))
                log("MHD_COMMENT is set");
            if (ar_tell(ar->stream) - ar->entry_offset > header.size) {
                warn("Invalid RAR header size: %d", header.size);
                return false;
            }
            rar->archive_flags = header.flags;
            break;

        case TYPE_FILE_ENTRY:
            if (!rar_parse_header_entry(rar, &header, &entry))
                return false;
            if ((header.flags & LHD_PASSWORD))
                warn("Encrypted entries will fail to uncompress");
            if ((header.flags & LHD_DIRECTORY) == LHD_DIRECTORY) {
                if (header.datasize == 0) {
                    log("Skipping directory entry \"%s\"", rar_get_name(ar, false));
                    break;
                }
                warn("Can't skip directory entries containing data");
            }
            if ((header.flags & (LHD_SPLIT_BEFORE | LHD_SPLIT_AFTER)))
                log("Splitting files across volumes");
            ar->entry_size_uncompressed = (size_t)entry.size;
            ar->entry_filetime = ar_conv_dosdate_to_filetime(entry.dosdate);
            if (!rar->entry.solid || rar->entry.method == METHOD_STORE || out_of_order) {
                rar_clear_uncompress(&rar->uncomp);
                memset(&rar->solid, 0, sizeof(rar->solid));
            }
            else {
                br_clear_leftover_bits(&rar->uncomp);
            }

            rar->solid.restart = rar->entry.solid && (out_of_order || !rar->solid.part_done);
            rar->solid.part_done = !ar->entry_size_uncompressed;
            rar->progress.data_left = (size_t)header.datasize;
            rar->progress.bytes_done = 0;
            rar->progress.crc = 0;

            /* TODO: CRC checks don't always hold (claim in XADRARParser.m @readBlockHeader) */
            if (!rar_check_header_crc(ar))
                warn("Invalid header checksum @%" PRIi64, ar->entry_offset);
            if (ar_tell(ar->stream) != ar->entry_offset + rar->entry.header_size) {
                warn("Couldn't seek to offset %" PRIi64, ar->entry_offset + rar->entry.header_size);
                return false;
            }
            return true;

        case TYPE_NEWSUB:
            log("Skipping newsub header @%" PRIi64, ar->entry_offset);
            break;

        case TYPE_END_OF_ARCHIVE:
            ar->at_eof = true;
            return false;

        default:
            log("Unknown RAR header type %02x", header.type);
            break;
        }

        /* TODO: CRC checks don't always hold (claim in XADRARParser.m @readBlockHeader) */
        if (!rar_check_header_crc(ar))
            warn("Invalid header checksum @%" PRIi64, ar->entry_offset);
        if (!ar_seek(ar->stream, ar->entry_offset_next, SEEK_SET)) {
            warn("Couldn't seek to offset %" PRIi64, ar->entry_offset_next);
            return false;
        }
    }
}

/* Continue reading from next volume for split files */
bool rar_continue_volume(ar_archive *ar)
{
    ar_archive_rar *rar = (ar_archive_rar *)ar;
    struct rar_header header;
    struct rar_entry entry;
    
    /* Check if we're still within the uncompressed data */
    if (rar->progress.bytes_done >= ar->entry_size_uncompressed)
        return false;
    
    /* The stream should now be positioned at the end of the current volume's data.
     * The multivolume stream will have automatically moved to the next volume.
     * Now we need to parse the header in the next volume to get the datasize. */
    
    /* Parse the header at current position (should be continuation header in next volume) */
    if (!rar_parse_header(ar, &header))
        return false;
    
    /* Check if this is a continuation header */
    if (header.type == TYPE_FILE_ENTRY) {
        /* Parse the entry details to get the header size */
        if (!rar_parse_header_entry(rar, &header, &entry))
            return false;
        
        /* Update data_left with the datasize from this volume part */
        rar->progress.data_left = (size_t)header.datasize;
        
        return true;
    } else if (header.type == TYPE_END_OF_ARCHIVE || header.type == TYPE_FILE_SIGNATURE || header.type == TYPE_MAIN_HEADER) {
        /* Skip non-file-entry headers and try again */
        if (!ar_skip(ar->stream, header.size - 7 + header.datasize))
            return false;
        
        /* Try parsing the next header recursively */
        return rar_continue_volume(ar);
    }
    
    return false;
}

static bool rar_copy_stored(ar_archive_rar *rar, void *buffer, size_t count)
{
    size_t total_read = 0;
    
    while (count > 0) {
        size_t to_read;
        
        /* If no data left in current volume part, try to continue to next volume */
        if (rar->progress.data_left == 0) {
            if (!rar_continue_volume(&rar->super))
                break;
        }
        
        to_read = count < rar->progress.data_left ? count : rar->progress.data_left;
        
        if (ar_read(rar->super.stream, (char *)buffer + total_read, to_read) != to_read) {
            warn("Unexpected EOF in stored data");
            return false;
        }
        
        rar->progress.data_left -= to_read;
        rar->progress.bytes_done += to_read;
        total_read += to_read;
        count -= to_read;
    }
    
    return total_read > 0 || count == 0;
}

static bool rar_restart_solid(ar_archive *ar)
{
    ar_archive_rar *rar = (ar_archive_rar *)ar;
    off64_t current_offset = ar->entry_offset;
    log("Restarting decompression for solid entry");
    if (!ar_parse_entry_at(ar, ar->entry_offset_first)) {
        ar_parse_entry_at(ar, current_offset);
        return false;
    }
    while (ar->entry_offset < current_offset) {
        size_t size = ar->entry_size_uncompressed;
        rar->solid.restart = false;
        while (size > 0) {
            unsigned char buffer[1024];
            size_t count = smin(size, sizeof(buffer));
            if (!ar_entry_uncompress(ar, buffer, count)) {
                ar_parse_entry_at(ar, current_offset);
                return false;
            }
            size -= count;
        }
        if (!ar_parse_entry(ar)) {
            ar_parse_entry_at(ar, current_offset);
            return false;
        }
    }
    rar->solid.restart = false;
    return true;
}

static bool rar_uncompress(ar_archive *ar, void *buffer, size_t count)
{
    ar_archive_rar *rar = (ar_archive_rar *)ar;
    if (count > ar->entry_size_uncompressed - rar->progress.bytes_done) {
        warn("Requesting too much data (%" PRIuPTR " < %" PRIuPTR ")", ar->entry_size_uncompressed - rar->progress.bytes_done, count);
        return false;
    }
    if (rar->entry.method == METHOD_STORE) {
        if (!rar_copy_stored(rar, buffer, count))
            return false;
    }
    else if (rar->entry.method == METHOD_FASTEST || rar->entry.method == METHOD_FAST ||
             rar->entry.method == METHOD_NORMAL || rar->entry.method == METHOD_GOOD ||
             rar->entry.method == METHOD_BEST) {
        if (rar->solid.restart && !rar_restart_solid(ar)) {
            warn("Failed to produce the required solid decompression state");
            return false;
        }
        if (!rar_uncompress_part(rar, buffer, count))
            return false;
    }
    else {
        warn("Unknown compression method %#02x", rar->entry.method);
        return false;
    }

    rar->progress.crc = ar_crc32(rar->progress.crc, buffer, count);
    if (rar->progress.bytes_done < ar->entry_size_uncompressed)
        return true;
    if (rar->progress.data_left)
        log("Compressed block has more data than required");
    rar->solid.part_done = true;
    rar->solid.size_total += rar->progress.bytes_done;
    if (rar->progress.crc != rar->entry.crc) {
        warn("Checksum of extracted data doesn't match");
        return false;
    }
    return true;
}

ar_archive *ar_open_rar_archive(ar_stream *stream)
{
    char signature[FILE_SIGNATURE_SIZE];
    if (!ar_seek(stream, 0, SEEK_SET))
        return NULL;
    if (ar_read(stream, signature, sizeof(signature)) != sizeof(signature))
        return NULL;
    if (memcmp(signature, "Rar!\x1A\x07\x00", sizeof(signature)) != 0) {
        if (memcmp(signature, "Rar!\x1A\x07\x01", sizeof(signature)) == 0)
            warn("RAR 5 format isn't supported");
        else if (memcmp(signature, "RE~^", 4) == 0)
            warn("Ancient RAR format isn't supported");
        else if (memcmp(signature, "MZ", 2) == 0 || memcmp(signature, "\x7F\x45LF", 4) == 0)
            warn("SFX archives aren't supported");
        return NULL;
    }

    return ar_open_archive(stream, sizeof(ar_archive_rar), rar_close, rar_parse_entry, rar_get_name, rar_uncompress, NULL, FILE_SIGNATURE_SIZE);
}
