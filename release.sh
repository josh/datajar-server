#!/bin/bash

set -euo pipefail
set -x

mkdir -p dist/
GOOS=darwin GOARCH=arm64 go build -o dist/datajar-server ./cmd/datajar-server
GOOS=linux GOARCH=amd64 go build -o dist/datajar-credential-server ./cmd/datajar-credential-server
