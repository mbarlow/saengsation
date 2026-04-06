#!/usr/bin/env bash
# Install saengsation into the uv environment so `saengsation` is available as a command.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_DIR"

REAL_HOME="$(getent passwd "${SUDO_USER:-$USER}" | cut -d: -f6)"
UV="${UV:-$(command -v uv || echo "$REAL_HOME/.local/bin/uv")}"
echo "Installing saengsation..."
"$UV" sync
"$UV" pip install -e .
echo
echo "Installed. Run with:"
echo "  uv run saengsation --help"
echo
echo "Or activate the venv and use directly:"
echo "  source .venv/bin/activate"
echo "  saengsation --help"
