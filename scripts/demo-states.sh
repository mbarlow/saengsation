#!/usr/bin/env bash
# Cycle through all built-in states with a pause between each.
# Usage: sudo ./scripts/demo-states.sh
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_DIR"

REAL_HOME="$(getent passwd "${SUDO_USER:-$USER}" | cut -d: -f6)"
UV="${UV:-$(command -v uv || echo "$REAL_HOME/.local/bin/uv")}"
S="$UV run saengsation"
PAUSE=4

echo "=== Saengsation State Demo ==="
echo "Cycling through all states (${PAUSE}s each)."
echo "Press Ctrl+C to stop."
echo

echo "Current status:"
$S status
echo

states=("focus" "night" "matrix" "chill" "music" "alert" "meeting")

for name in "${states[@]}"; do
    echo "-----------------------------"
    echo ">>> state: $name"
    $S state show "$name" | head -1
    $S state set "$name"
    echo "    (watching for ${PAUSE}s...)"
    sleep "$PAUSE"
    echo
done

echo "-----------------------------"
echo ">>> Restoring: focus"
$S state set focus

echo
echo "Demo complete!"
