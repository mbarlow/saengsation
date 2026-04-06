#!/usr/bin/env bash
# Check that all dependencies and system requirements are met.
set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

ok()   { echo -e "  ${GREEN}[OK]${NC}  $1"; }
warn() { echo -e "  ${YELLOW}[!!]${NC}  $1"; }
fail() { echo -e "  ${RED}[NO]${NC}  $1"; }

errors=0

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
BINARY="$PROJECT_DIR/saengsation"

echo "Checking saengsation setup..."
echo

# Binary
echo "Checking binary..."
if [ -x "$BINARY" ]; then
    ok "saengsation binary found"
else
    fail "saengsation binary not found — run: make build"
    ((errors++))
fi

# udev rules
echo
echo "Checking udev rules..."
if [ -f /etc/udev/rules.d/99-saengsation.rules ]; then
    ok "udev rules installed"
else
    warn "udev rules not installed — run: make setup-udev"
fi

# plugdev group
echo
echo "Checking permissions..."
if getent group plugdev &>/dev/null; then
    ok "plugdev group exists"
    if id -nG 2>/dev/null | grep -qw plugdev; then
        ok "current user is in plugdev"
    else
        warn "current user is NOT in plugdev — run: make setup-group"
    fi
else
    warn "plugdev group does not exist — run: make setup-group"
fi

# Devices
echo
echo "Checking devices..."
if lsusb 2>/dev/null | grep -q "3434:0370"; then
    ok "Keychron V7 detected"
else
    warn "Keychron V7 not detected (not plugged in?)"
fi

# hidraw permissions
echo
echo "Checking hidraw access..."
for dev in /dev/hidraw*; do
    if [ -w "$dev" ] 2>/dev/null; then
        product=$(udevadm info --query=all --name="$dev" 2>/dev/null | grep ID_MODEL= | head -1 | cut -d= -f2)
        if [ -n "$product" ]; then
            ok "$dev writable ($product)"
        fi
    fi
done

echo
if [ "$errors" -gt 0 ]; then
    echo -e "${RED}$errors issue(s) found.${NC}"
    exit 1
else
    echo -e "${GREEN}All checks passed.${NC}"
fi
