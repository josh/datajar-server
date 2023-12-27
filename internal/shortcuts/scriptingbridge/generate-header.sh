#!/bin/bash

set -euo pipefail
set -x

sdef /System/Applications/Shortcuts.app | sdp -fh --basename Shortcuts
