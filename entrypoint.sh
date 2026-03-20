#!/usr/bin/env bash

set -euo pipefail

#fallback_to_shell() {
#  echo "Trapped SIGINT. Falling back to /bin/bash..." 1>&2
#  exec /bin/bash
#}
#
#trap fallback_to_shell INT

set +u
if [ -n "$1" ]; then
  echo "Overriding to run: ${@:1}"
  exec ${@:1}
else
  set -u
  echo 'About to spawn lighttpd...'

  exec lighttpd -D -f /etc/lighttpd/lighttpd.conf
fi

echo 'Unreachable!' 1>&2
exit 1
