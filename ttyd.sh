#!/usr/bin/env bash

set -euo pipefail

ROOT="$(dirname "$(realpath "${BASH_SOURCE[0]}")" )"
PID_FILE="${ROOT}/ttyd.pid"

cd "$ROOT"

# --credential user:pw
ttyd --writable --cwd "$HOME" /bin/bash &
PID=$!

echo "Forked ttyd as the background process $PID"
echo "$!" > "$PID_FILE"

wait $!

rm "$PID_FILE"
