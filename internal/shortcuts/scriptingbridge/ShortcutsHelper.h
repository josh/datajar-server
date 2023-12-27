typedef struct
{
    const void *bytes;
    int length;
    const char *err;
} ShortcutResult;

void runShortcut(const char *shortcutName, ShortcutResult *result);
