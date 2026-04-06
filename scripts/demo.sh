#!/usr/bin/env bash
# Run a sequence of demo animations to show off saengsation.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
S="$PROJECT_DIR/saengsation"

run() {
    echo
    echo ">>> saengsation $*"
    "$S" "$@"
}

echo "=== Saengsation Demo ==="
echo "This will cycle through keyboard animations."
echo "Press Ctrl+C at any time to stop."
echo

# Show what's connected
run status

sleep 1

echo
echo "--- Color Cycle (5s) ---"
run animate cycle --duration 5

sleep 0.5

echo
echo "--- Police Lights (5s) ---"
run animate police --duration 5

sleep 0.5

echo
echo "--- Blue Pulse (5s) ---"
run animate pulse --hue 170 --duration 5

sleep 0.5

echo
echo "--- Red Flash (3s) ---"
run animate flash --hue 0 --duration 3

sleep 0.5

echo
echo "--- Green Flash (3s) ---"
run animate flash --hue 85 --duration 3

sleep 0.5

# Reset to a calm state
echo
echo "--- Resetting to breathing effect ---"
"$S" kb effect breathing 2>/dev/null || true

echo
echo "Demo complete!"
