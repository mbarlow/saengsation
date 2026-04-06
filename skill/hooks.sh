#!/usr/bin/env bash
# Saengsation hooks for Claude Code.
# Changes keyboard lighting based on Claude Code events.
#
# Receives JSON on stdin from Claude Code (ignored — we only use the argument).
# Runs saengsation in the background so it doesn't block Claude.
set -euo pipefail

# Drain stdin so Claude Code doesn't hang
cat > /dev/null &

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
SAENGSATION="${SAENGSATION:-$PROJECT_DIR/saengsation}"

STATE="${1:-idle}"

"$SAENGSATION" state set "$STATE" 2>/dev/null &

exit 0
