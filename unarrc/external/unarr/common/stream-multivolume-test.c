/* Test helper for multi-volume naming */
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

/* Get the next volume path based on current volume number
 * Naming convention:
 *   -1: .rar
 *   0-99: .r00 to .r99
 *   100-199: .s00 to .s99
 *   200-299: .t00 to .t99
 *   300-399: .u00 to .u99
 *   400-499: .v00 to .v99
 *   500-599: .w00 to .w99
 *   600-699: .x00 to .x99
 *   700-799: .y00 to .y99
 *   800-899: .z00 to .z99 (900 total volumes)
 */
char *test_get_volume_path(const char *base_path, long long volume_num) {
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
        
        /* Support volumes 0-899 (.r00 through .z99) */
        if (letter > 'z') {
            free(path);
            return NULL; /* Volume number too high */
        }
        
        sprintf(path + base_len, ".%c%02d", letter, remainder);
    }
    
    return path;
}

void test_free_path(char *path) {
    free(path);
}
