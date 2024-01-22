typedef struct
{
    const void *bytes;
    int length;
    const char *err;
} ShortcutResult;

int hasShortcut(const char *name);
void runShortcut(const char *name, const char *input, ShortcutResult *result);
