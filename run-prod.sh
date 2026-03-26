#!/usr/bin/env bash

# In bash since we don't get have a slothscript prebuilt for armv7

# Don't set -u since we read CLI args
set -eo pipefail

ROOT="$(dirname "$(realpath "${BASH_SOURCE[0]}" )" )"

ROS_PROPRIETARY="$(realpath "${ROOT}/../")"

LOCAL_REPOS="$1"

if [[ ! -d "${LOCAL_REPOS}" ]]; then
  if [[ -z "$1" ]]; then
    echo "You must pass a PATH_TO_REPOS" >&2
    echo
  else
    echo "No local repos at \"${LOCAL_REPOS}\"" >&2
    echo
  fi
	echo "Usage: run-prod.sh PATH_TO_REPOS [DEBUG_OVERRIDE_CMD ...]" >&2
	exit 1
fi

cd "$ROOT"

if [[ -n "$2" ]]; then
  CMD=${@:2}
  echo "debug mode with command: ${CMD}"

  # Don't wrap $CMD in quotes
  docker container run \
    -it \
    --rm \
    --name ros \
    --mount type=bind,src="${LOCAL_REPOS}",dst=/repos \
    --mount type=bind,src="${ROS_PROPRIETARY}",dst=/ros-proprietary \
    --publish 80:80 \
    cgit-env \
    $CMD
else
  echo "prod mode"

  # Without --detach, stdout & stderr will stream to term and command will block
  # until the container finishes.
  docker container run \
    --rm \
    --detach \
    --name ros \
    --mount type=bind,src="${LOCAL_REPOS}",dst=/repos \
    --mount type=bind,src="${ROS_PROPRIETARY}",dst=/ros-proprietary \
    --publish 80:80 \
    cgit-env
fi
