#!/usr/bin/env bash
# Install saengsation hooks into Claude Code settings.
# Merges hook entries into ~/.claude/settings.json without overwriting existing settings.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
HOOKS_SH="$PROJECT_DIR/skill/hooks.sh"
SETTINGS_FILE="$HOME/.claude/settings.json"

if [ ! -f "$HOOKS_SH" ]; then
    echo "ERROR: hooks.sh not found at $HOOKS_SH"
    exit 1
fi

mkdir -p "$HOME/.claude"

# Build the hooks JSON with the correct absolute path
HOOKS_JSON=$(cat <<EOF
{
  "UserPromptSubmit": [
    {
      "matcher": "",
      "hooks": [
        {
          "type": "command",
          "command": "$HOOKS_SH acknowledged"
        }
      ]
    }
  ],
  "PreToolUse": [
    {
      "matcher": "",
      "hooks": [
        {
          "type": "command",
          "command": "$HOOKS_SH working"
        }
      ]
    }
  ],
  "Stop": [
    {
      "matcher": "",
      "hooks": [
        {
          "type": "command",
          "command": "$HOOKS_SH idle"
        }
      ]
    }
  ],
  "Notification": [
    {
      "matcher": "idle_prompt",
      "hooks": [
        {
          "type": "command",
          "command": "$HOOKS_SH waiting"
        }
      ]
    }
  ]
}
EOF
)

# Merge into existing settings or create new
if [ -f "$SETTINGS_FILE" ]; then
    # Check if jq is available for safe merging
    if command -v jq &>/dev/null; then
        EXISTING=$(cat "$SETTINGS_FILE")
        echo "$EXISTING" | jq --argjson hooks "$HOOKS_JSON" '.hooks = ($hooks * (.hooks // {}))' > "$SETTINGS_FILE.tmp"
        mv "$SETTINGS_FILE.tmp" "$SETTINGS_FILE"
        echo "Merged saengsation hooks into $SETTINGS_FILE"
    else
        echo "WARNING: jq not installed — cannot safely merge into existing settings."
        echo ""
        echo "Add the following hooks to $SETTINGS_FILE manually:"
        echo ""
        echo "  \"hooks\": $(echo "$HOOKS_JSON" | head -1)"
        echo "  ..."
        echo ""
        echo "Or install jq and re-run:  sudo pacman -S jq  (or apt install jq)"
        exit 1
    fi
else
    # No existing settings — create fresh
    if command -v jq &>/dev/null; then
        jq -n --argjson hooks "$HOOKS_JSON" '{"hooks": $hooks}' > "$SETTINGS_FILE"
    else
        cat > "$SETTINGS_FILE" <<SETTINGS
{
  "hooks": $HOOKS_JSON
}
SETTINGS
    fi
    echo "Created $SETTINGS_FILE with saengsation hooks."
fi

echo ""
echo "Hooks installed. Claude Code will now change your keyboard lighting:"
echo "  UserPromptSubmit → acknowledged (green breathing)"
echo "  PreToolUse       → working (rainbow spiral)"
echo "  Stop             → idle (dim blue pulse)"
echo "  Notification     → waiting (solid red)"
echo ""
echo "Restart Claude Code for hooks to take effect."
