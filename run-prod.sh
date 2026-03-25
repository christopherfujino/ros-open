#!/usr/bin/env bash

# In bash since we don't get have a slothscript prebuilt for armv7

set -euo pipefail

ROOT="$(dirname "$(realpath "${BASH_SOURCE[0]}" )" )"

set +u
LOCAL_REPOS="$1"
set -u

if [[ ! -d "${LOCAL_REPOS}" ]]; then
	echo "No local repos at \"${LOCAL_REPOS}\"" >&2
	echo "Usage: run-prod.sh [PATH TO REPOS]" >&2
	exit 1
fi

cd "$ROOT"

# Without --detach, stdout & stderr will stream to term and command will block
# until the container finishes.
docker container run \
	--rm \
	--detach \
	--name cgit-env \
	--mount type=bind,src="$LOCAL_REPOS",dst=/repos \
	--publish 80:80 \
	cgit-env
