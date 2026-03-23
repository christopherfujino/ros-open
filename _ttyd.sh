# vim: ft=bash

set -euo pipefail

set +u
if [[ -z "$ENSURE_CALLED_FROM_WRAPPER" ]]; then
	echo "Do not invoke this script directly, instead call `./ttyd.sh`" >&2
	exit 1
fi
set -u

ROOT="$(dirname "$(realpath "${BASH_SOURCE[0]}")" )"
PID_FILE="${ROOT}/ttyd.pid"

cd "$ROOT"

# --credential user:pw
ttyd --writable --cwd "$HOME" /bin/bash &
PID=$!

echo "$!" > "$PID_FILE"

sleep 0.375 # let ttyd output flush first...
echo "Forked ttyd as the background process $PID"
echo 'To check ttyd:'
echo '    ps $(cat ./ttyd.pid)'
echo
echo 'To kill ttyd:'
echo '    kill $(cat ./ttyd.pid)'
echo
echo 'Press enter to regain shell access'
wait $!

rm "$PID_FILE"
echo "Removed $PID_FILE"
