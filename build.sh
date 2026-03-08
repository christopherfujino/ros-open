#!/usr/bin/env bash

set -euo pipefail

ROOT="$(dirname "$(realpath "${BASH_SOURCE[0]}" )" )"
docker build "$ROOT" --tag cgit-env
