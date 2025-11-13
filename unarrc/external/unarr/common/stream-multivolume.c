/* Multi-volume archive stream support for RAR archives */
/* Copyright 2025 - License: LGPLv3 */

#include "unarr-imp.h"
#include <string.h>
#include <stdio.h>

struct MultiVolumeStream {
    char *base_path;          /* Base path without extension */
    char *current_path;       /* Current volume path */
    FILE *current_file;       /* Current open file handle */
    off64_t current_volume;   /* Current volume number (-1 for .rar, 0+ for .r00, .r01, ..., .z99) */
    off64_t offset;           /* Global offset across all volumes */
    off64_t current_size;     /* Size of current volume */
    off64_t total_offset;     /* Total bytes read so far */
};

/* Get the next volume path based on current volume number
 * Naming convention:
 *   -1: .rar
 *   0-99: .r00 to .r99
 *   100-199: .s00 to .s99
 *   200-299: .t00 to .t99
 *   300-399: .u00 to .u99
 *   ... continuing through alphabet to .z99 (2599 total volumes)
 */
static char *get_volume_path(const char *base_path, off64_t volume_num) {
    size_t base_len = strlen(base_path);
    char *path = malloc(base_len + 10); /* Room for extension */
    if (!path)
        return NULL;
    
    strcpy(path, base_path);
    
    if (volume_num < 0) {
        /* First volume */
        strcat(path, ".rar");
    } else {
        /* Calculate the letter suffix: r, s, t, u, ..., z */
        int hundred = (int)(volume_num / 100);
        int remainder = (int)(volume_num % 100);
        char letter = 'r' + hundred;
        
        /* Support volumes 0-2599 (.r00 through .z99) */
        if (letter > 'z') {
            free(path);
            return NULL; /* Volume number too high */
        }
        
        sprintf(path + base_len, ".%c%02d", letter, remainder);
    }
    
    return path;
}

/* Get file size */
static off64_t get_file_size(FILE *f) {
    off64_t current_pos, size;
    
    current_pos = ftell(f);
    if (current_pos < 0)
        return -1;
    
    if (fseek(f, 0, SEEK_END) != 0)
        return -1;
    
    size = ftell(f);
    
    if (fseek(f, current_pos, SEEK_SET) != 0)
        return -1;
    
    return size;
}

/* Try to open the next volume */
static bool multivolume_switch_to_volume(struct MultiVolumeStream *mv, off64_t volume_num) {
    FILE *new_file;
    char *new_path;
    
    /* Already on this volume? */
    if (volume_num == mv->current_volume && mv->current_file != NULL)
        return true;
    
    /* Generate path for the new volume */
    new_path = get_volume_path(mv->base_path, volume_num);
    if (!new_path)
        return false;
    
    /* Try to open the volume */
    new_file = fopen(new_path, "rb");
    if (!new_file) {
        free(new_path);
        return false;
    }
    
    /* Close the old file if open */
    if (mv->current_file) {
        fclose(mv->current_file);
    }
    
    /* Update state */
    free(mv->current_path);
    mv->current_path = new_path;
    mv->current_file = new_file;
    mv->current_volume = volume_num;
    mv->current_size = get_file_size(new_file);
    
    return true;
}

static void multivolume_close(void *data) {
    struct MultiVolumeStream *mv = data;
    if (mv->current_file)
        fclose(mv->current_file);
    free(mv->base_path);
    free(mv->current_path);
    free(mv);
}

static size_t multivolume_read(void *data, void *buffer, size_t count) {
    struct MultiVolumeStream *mv = data;
    size_t total_read = 0;
    
    if (!mv->current_file)
        return 0;
    
    while (count > 0) {
        size_t bytes_read;
        off64_t current_pos = ftell(mv->current_file);
        size_t available;
        
        if (current_pos < 0)
            break;
        
        /* How much is available in this volume? */
        available = (size_t)(mv->current_size - current_pos);
        
        if (available == 0) {
            /* Need to switch to next volume */
            if (!multivolume_switch_to_volume(mv, mv->current_volume + 1))
                break; /* No more volumes */
            
            /* Seek to start of new volume */
            if (fseek(mv->current_file, 0, SEEK_SET) != 0)
                break;
            
            continue;
        }
        
        /* Read from current volume */
        bytes_read = fread((char *)buffer + total_read, 1, 
                          (count < available) ? count : available, 
                          mv->current_file);
        
        if (bytes_read == 0)
            break;
        
        total_read += bytes_read;
        count -= bytes_read;
        mv->total_offset += bytes_read;
    }
    
    return total_read;
}

static bool multivolume_seek(void *data, off64_t offset, int origin) {
    struct MultiVolumeStream *mv = data;
    off64_t target_offset;
    off64_t cumulative_offset = 0;
    off64_t volume_num = -1;  /* Start with .rar */
    
    if (!mv->current_file)
        return false;
    
    /* Calculate absolute target offset */
    if (origin == SEEK_SET) {
        target_offset = offset;
    } else if (origin == SEEK_CUR) {
        target_offset = mv->total_offset + offset;
    } else {
        /* SEEK_END not supported for multi-volume */
        return false;
    }
    
    if (target_offset < 0)
        return false;
    
    /* Find which volume contains this offset */
    /* This is a simplified approach - we try volumes in sequence */
    /* A better implementation would cache volume sizes */
    
    while (true) {
        off64_t volume_size;
        
        if (!multivolume_switch_to_volume(mv, volume_num))
            return false;
        
        volume_size = mv->current_size;
        
        if (target_offset < cumulative_offset + volume_size) {
            /* This volume contains our target offset */
            off64_t local_offset = target_offset - cumulative_offset;
            
            if (fseek(mv->current_file, (long)local_offset, SEEK_SET) != 0)
                return false;
            
            mv->total_offset = target_offset;
            return true;
        }
        
        cumulative_offset += volume_size;
        volume_num++;
        
        if (volume_num > 2599) /* Safety limit: max .z99 */
            return false;
    }
}

static off64_t multivolume_tell(void *data) {
    struct MultiVolumeStream *mv = data;
    return mv->total_offset;
}

/* Open a multi-volume RAR archive */
ar_stream *ar_open_multivolume(const char *path) {
    struct MultiVolumeStream *mv;
    char *base_path;
    const char *ext;
    size_t base_len;
    
    if (!path)
        return NULL;
    
    /* Find the extension */
    ext = strrchr(path, '.');
    if (!ext)
        return NULL;
    
    /* Extract base path (without extension) */
    base_len = ext - path;
    base_path = malloc(base_len + 1);
    if (!base_path)
        return NULL;
    
    memcpy(base_path, path, base_len);
    base_path[base_len] = '\0';
    
    /* Allocate multi-volume state */
    mv = calloc(1, sizeof(struct MultiVolumeStream));
    if (!mv) {
        free(base_path);
        return NULL;
    }
    
    mv->base_path = base_path;
    mv->current_volume = -2; /* Not opened yet */
    mv->offset = 0;
    mv->total_offset = 0;
    
    /* Open first volume (.rar) */
    if (!multivolume_switch_to_volume(mv, -1)) {
        free(base_path);
        free(mv);
        return NULL;
    }
    
    return ar_open_stream(mv, multivolume_close, multivolume_read, multivolume_seek, multivolume_tell);
}
