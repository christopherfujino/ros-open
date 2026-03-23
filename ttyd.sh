#!/usr/bin/env bash

set -euo pipefail

ROOT="$(dirname "$(realpath "${BASH_SOURCE[0]}")" )"

ENSURE_CALLED_FROM_WRAPPER=1 bash "${ROOT}/_ttyd.sh" &
