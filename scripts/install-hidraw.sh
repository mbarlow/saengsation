#!/usr/bin/env bash
# Build and install the hidapi Python module with the hidraw backend.
#
# The PyPI hidapi wheel bundles libusb, which requires root to detach the
# kernel driver.  The hidraw backend uses /dev/hidraw* devices that can be
# permissioned via udev rules — no root needed at runtime.
#
# Requires: git, pkg-config, gcc, system hidapi library (e.g. pacman -S hidapi)
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_DIR"

REAL_HOME="$(getent passwd "${SUDO_USER:-$USER}" | cut -d: -f6)"
UV="${UV:-$(command -v uv 2>/dev/null || echo "$REAL_HOME/.local/bin/uv")}"
PYTHON="$PROJECT_DIR/.venv/bin/python"

# Verify system hidapi-hidraw is available
if ! pkg-config --exists hidapi-hidraw 2>/dev/null; then
    echo "ERROR: hidapi-hidraw not found via pkg-config."
    echo "Install the system hidapi library first:"
    echo "  Arch:   sudo pacman -S hidapi"
    echo "  Debian: sudo apt install libhidapi-hidraw0 libhidapi-dev"
    echo "  Fedora: sudo dnf install hidapi-devel"
    exit 1
fi

# Ensure venv exists with build deps
"$UV" sync
"$UV" pip install cython setuptools

# Clone cython-hidapi to a temp directory
TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT

echo "Cloning cython-hidapi..."
git clone --depth 1 --quiet https://github.com/trezor/cython-hidapi.git "$TMPDIR/cython-hidapi"

cd "$TMPDIR/cython-hidapi"

# The .pyx source is named hidraw.pyx but we need the module to export as
# "hid" so that `import hid` works.  Cythonize with the correct module name.
echo "Building hidapi with hidraw backend..."
"$PYTHON" -m cython hid.pyx -o hid.c
"$PYTHON" -c "
from setuptools import setup, Extension
import subprocess

cflags = subprocess.check_output(['pkg-config', '--cflags', 'hidapi-hidraw']).decode().strip().split()
libs = subprocess.check_output(['pkg-config', '--libs', 'hidapi-hidraw']).decode().strip().split()

ext = Extension('hid', sources=['hid.c'])
ext.extra_compile_args = cflags
ext.extra_link_args = libs

setup(
    name='hidapi',
    ext_modules=[ext],
    script_args=['build_ext', '--inplace']
)
"

# Remove any existing hidapi pip artifacts and install our build
SITE_PACKAGES="$("$PYTHON" -c 'import site; print(site.getsitepackages()[0])')"
rm -f "$SITE_PACKAGES"/hid.cpython-*.so
rm -rf "$SITE_PACKAGES"/hidapi.libs
rm -rf "$SITE_PACKAGES"/hidapi-*.dist-info
rm -rf "$SITE_PACKAGES"/hidapi-*.egg-info
rm -f "$SITE_PACKAGES"/hidraw.cpython-*.so

cp hid.cpython-*.so "$SITE_PACKAGES/"

echo "Verifying hidraw backend..."
if "$PYTHON" -c "import hid; hid.enumerate()" 2>/dev/null; then
    echo "  hidapi installed with hidraw backend."
else
    echo "  WARNING: hidapi installed but enumerate failed. Check system hidapi library."
fi

echo "Done."
