#!/usr/bin/env zsh
# Helper script to run the server with CGO disabled
set -euo pipefail

export CGO_ENABLED=0
exec go run ./cmd
