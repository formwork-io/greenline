#/usr/bin/env bash

# The next three lines are for the go shell.
export SCRIPT_NAME="test-loop"
export SCRIPT_HELP="Run make test when src files are modified."
[[ "$GOGO_GOSH_SOURCE" -eq 1 ]] && return 0

# Normal script execution starts here.
dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"/../
source "$dir"/env.sh || exit 1
assert-env-or-die src
cd "$dir" || exit 1

while true; do
    rslt=$("$GOSH_SCRIPTS"/gate $(find "$src" -type d -print)) || exit 1
    clear
    dir_changed="$(dirname $(echo $rslt | awk '{ print $1 }'))"
    date +"%a %b %m %T"
    echo "[TESTING]"
    cd "$dir_changed" || exit 1
    go test
    cd -
done
exit 0

