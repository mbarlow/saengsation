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
REAL_HOME="$(getent passwd "${SUDO_USER:-$USER}" | cut -d: -f6)"
UV="${UV:-$(command -v uv 2>/dev/null || echo "$REAL_HOME/.local/bin/uv")}"

STATE="${1:-idle}"

cd "$PROJECT_DIR"
"$UV" run saengsation state set "$STATE" 2>/dev/null &

exit 0
