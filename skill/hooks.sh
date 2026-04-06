#!/usr/bin/env bash
# Saengsation hooks for Claude Code.
# These change keyboard lighting based on Claude Code events.
#
# Install by adding to your Claude Code settings (~/.claude/settings.json):
#
#   "hooks": {
#     "PreToolUse": [
#       { "command": "/path/to/saengsation/skill/hooks.sh working" }
#     ],
#     "PostToolUse": [
#       { "command": "/path/to/saengsation/skill/hooks.sh working" }
#     ],
#     "Notification": [
#       { "command": "/path/to/saengsation/skill/hooks.sh waiting" }
#     ],
#     "Stop": [
#       { "command": "/path/to/saengsation/skill/hooks.sh idle" }
#     ]
#   }
#
# Or use the granular version below for full lifecycle control.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
REAL_HOME="$(getent passwd "${SUDO_USER:-$USER}" | cut -d: -f6)"
UV="${UV:-$(command -v uv 2>/dev/null || echo "$REAL_HOME/.local/bin/uv")}"

run() {
    cd "$PROJECT_DIR"
    "$UV" run saengsation state set "$1" 2>/dev/null &
}

case "${1:-}" in
    waiting)      run waiting ;;       # Claude needs your input
    acknowledged) run acknowledged ;;  # Claude received input
    working)      run working ;;       # Claude is processing
    idle)         run idle ;;          # Nothing happening
    *)            run "${1:-idle}" ;;  # Any custom state name
esac
