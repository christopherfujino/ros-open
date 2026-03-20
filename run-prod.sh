#!/usr/bin/env bash

# In bash since we don't get have a slothscript prebuilt for armv7

set -euo pipefail

ROOT="$(dirname "$(realpath "${BASH_SOURCE[0]}" )" )"
LOCAL_REPOS="${HOME}/repos"

cd "$ROOT"

# Without --detach, stdout & stderr will stream to term and command will block
# until the container finishes.
docker container run \
	--rm \
	--detach \
	--name cgit-env \
	--mount type=bind,src=./repos,dst=/repos \
	--publish 80:80 \
	cgit-env
