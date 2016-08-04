#!/usr/bin/env bash

# The next three lines are for the go shell.
export SCRIPT_NAME="make-loop"
export SCRIPT_HELP="Run make when src files are modified."
[[ "$GOGO_GOSH_SOURCE" -eq 1 ]] && return 0

# Normal script execution starts here.
dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"/../
source "$dir"/env.sh || exit 1
assert-env-or-die src
cd "$dir" || exit 1

while true; do
    "$GOSH_SCRIPTS"/gate $(find "$src" -type d -print) || exit 1
    clear
    date +"%a %b %m %T"
    echo "[BUILDING]"
    make
done
exit 0

