#!/usr/bin/env bash
# Normal script execution starts here.
dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"/../
source "$dir"/env.sh || exit 1
assert-env-or-die src
cd "$dir" || exit 1

while true; do
    "$GOSH_SCRIPTS"/gate $(find "$src" -type d -print) || exit 1
    clear
    date +"%a %b %m %T"
    "$@"
    echo "Completed @ $(date +"%a %b %m %T")"
done
exit 0

