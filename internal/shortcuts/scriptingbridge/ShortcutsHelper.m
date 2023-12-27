#import <Foundation/Foundation.h>
#import <ScriptingBridge/ScriptingBridge.h>
#import "Shortcuts.h"
#import "ShortcutsHelper.h"

void runShortcut(const char* name, const char *input, ShortcutResult *result) {
    @try {
        ShortcutsApplication *app = [SBApplication applicationWithBundleIdentifier:@"com.apple.shortcuts.events"];
        SBElementArray<ShortcutsShortcut *> *shortcuts = [app shortcuts];
        ShortcutsShortcut *shortcut = [shortcuts objectWithName:[NSString stringWithUTF8String:name]];

        if ([shortcut name] == nil) {
            result->err = "Shortcut not found";
            return;
        }

        id output = [shortcut runWithInput:[NSString stringWithUTF8String:input]];

        if (output == nil) {
            return;
        }

        NSError *writeError = nil;
        NSData *jsonData = [NSJSONSerialization dataWithJSONObject:output options:0 error:&writeError];

        if (writeError != nil) {
            result->err = [[writeError localizedDescription] UTF8String];
        } else {
            result->bytes = [jsonData bytes];
            result->length = [jsonData length];
        }
    } @catch (NSException *exception) {
        result->err = [[exception reason] UTF8String];
    }
}
