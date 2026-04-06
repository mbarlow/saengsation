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

echo "Checking saengsation dependencies..."
echo

# uv
if command -v uv &>/dev/null; then
    ok "uv $(uv --version 2>/dev/null | head -1)"
else
    fail "uv not found — install from https://docs.astral.sh/uv/"
    ((errors++))
fi

# Python
if command -v python3 &>/dev/null; then
    ok "python3 $(python3 --version 2>&1 | awk '{print $2}')"
else
    fail "python3 not found"
    ((errors++))
fi

# hidapi system library
if ldconfig -p 2>/dev/null | grep -q libhidapi || pacman -Qi hidapi &>/dev/null 2>&1; then
    ok "hidapi system library"
else
    fail "hidapi not found — install: sudo pacman -S hidapi"
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

# Python hidapi backend
echo
echo "Checking Python hidapi backend..."
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
PYTHON="$PROJECT_DIR/.venv/bin/python"
if [ -f "$PYTHON" ]; then
    HID_SO="$("$PYTHON" -c 'import importlib.util; print(importlib.util.find_spec("hid").origin)' 2>/dev/null || true)"
    if [ -n "$HID_SO" ]; then
        if ldd "$HID_SO" 2>/dev/null | grep -q libhidapi-hidraw; then
            ok "hidapi using hidraw backend"
        elif ldd "$HID_SO" 2>/dev/null | grep -q libusb; then
            fail "hidapi using libusb backend (needs root) — run: make install"
            ((errors++))
        else
            warn "could not determine hidapi backend"
        fi
    else
        fail "Python hid module not installed — run: make install"
        ((errors++))
    fi
else
    warn "venv not found — run: make install"
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
    echo -e "${RED}$errors required dependency missing.${NC}"
    exit 1
else
    echo -e "${GREEN}All required dependencies found.${NC}"
fi
