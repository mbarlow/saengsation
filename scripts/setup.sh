#!/usr/bin/env bash
# Full setup: build binary, create group, add user, install udev rules.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
RULES_FILE="$PROJECT_DIR/99-saengsation.rules"

echo "=== Saengsation Setup ==="
echo

# Build
echo "Building saengsation..."
cd "$PROJECT_DIR"
go build -o saengsation .
echo "  Binary built: $PROJECT_DIR/saengsation"

# plugdev group
if ! getent group plugdev &>/dev/null; then
    echo "Creating plugdev group..."
    sudo groupadd plugdev
fi

if ! id -nG | grep -qw plugdev; then
    echo "Adding $(whoami) to plugdev..."
    sudo usermod -aG plugdev "$(whoami)"
    echo "  NOTE: You must log out and back in (or run 'newgrp plugdev') for group changes to take effect."
fi

# udev rules
echo "Installing udev rules..."
sudo cp "$RULES_FILE" /etc/udev/rules.d/99-saengsation.rules
sudo udevadm control --reload-rules
sudo udevadm trigger
echo "  Udev rules installed and reloaded."

echo
echo "Setup complete. If you were added to plugdev, log out/in or run:"
echo "  newgrp plugdev"
echo
echo "Then test with:"
echo "  ./saengsation status"
