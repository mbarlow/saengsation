# saengsation

RGB LED controller for Keychron V7 on Linux, with Claude Code integration.

*แสง (saeng) = light + sensation*

Control your keyboard lighting from the command line, define named states for different activities, and let Claude Code change your keyboard colors based on what it's doing.

## Quickstart

```bash
# Clone
git clone git@github.com:mbarlow/saengsation.git
cd saengsation

# Setup (creates plugdev group, installs udev rules, builds hidraw backend)
make setup

# Log out and back in for group permissions, or:
newgrp plugdev

# Verify
make check

# Try it
uv run saengsation state set focus

# Optional: install Claude Code hooks
make hooks
```

## Install

### Requirements

- Linux (tested on Arch)
- Python 3.12+
- [uv](https://docs.astral.sh/uv/) package manager
- `hidapi` system library with hidraw support:
  - Arch: `sudo pacman -S hidapi`
  - Debian/Ubuntu: `sudo apt install libhidapi-hidraw0 libhidapi-dev`
  - Fedora: `sudo dnf install hidapi-devel`
- Build tools: `gcc`, `pkg-config`, `git` (for building the hidraw Python module)
- Keychron V7 keyboard (USB, QMK firmware)

### Step by Step

```bash
# Install Python deps and build hidraw backend
make install

# Create plugdev group and add yourself
sudo groupadd plugdev
sudo usermod -aG plugdev $USER

# Install udev rules (allows non-root HID access)
sudo cp 99-saengsation.rules /etc/udev/rules.d/
sudo udevadm control --reload-rules
sudo udevadm trigger

# Log out/in for group to take effect
```

Or just run `make setup` which does all of the above.

> **Why hidraw?** The `hidapi` pip package ships with a libusb backend that
> requires root to detach the kernel driver. `make install` builds it from
> source against the system `libhidapi-hidraw` library, which uses
> `/dev/hidraw*` devices permissioned via udev rules — no root needed.

## Usage

All commands below use `uv run` to run within the project environment. If you activate the venv (`source .venv/bin/activate`), you can omit the `uv run` prefix.

### Direct Control

```bash
# Show keyboard status
uv run saengsation status

# Set effect
uv run saengsation kb effect breathing
uv run saengsation kb effect digital_rain
uv run saengsation kb effect 5

# Set color (hue, saturation, brightness — each 0-255)
uv run saengsation kb color 85,255,200       # green
uv run saengsation kb color 0,255,255        # red
uv run saengsation kb color 170,255          # blue (keeps current brightness)

# Brightness and speed
uv run saengsation kb brightness 128
uv run saengsation kb speed 2                # 0-3

# Save to EEPROM (persists across unplugs)
uv run saengsation kb effect breathing --save

# List all effects
uv run saengsation effects
```

### Animations

```bash
uv run saengsation animate cycle                # rainbow color cycle
uv run saengsation animate police               # red/blue flash
uv run saengsation animate pulse --hue 170      # blue pulse
uv run saengsation animate flash --hue 0        # red strobe
uv run saengsation animate cycle --duration 30  # 30 seconds
```

### Named States

States are presets that bundle an effect, color, brightness, and speed into a single name.

```bash
# List all states
uv run saengsation state list

# Apply a state
uv run saengsation state set focus
uv run saengsation state set matrix

# Apply and save to EEPROM
uv run saengsation state set night --save

# View state details
uv run saengsation state show focus

# Save current keyboard settings as a new state
uv run saengsation state save mystate -d "Purple haze for coding"

# Delete a custom state
uv run saengsation state delete mystate
```

### Built-in States

| State | Effect | Description |
|-------|--------|-------------|
| `focus` | breathing blue | Deep work mode |
| `alert` | solid red | Something needs attention |
| `chill` | cycle_all | Slow rainbow cycle |
| `meeting` | off | Lights off |
| `music` | rainbow_beacon | Party mode |
| `night` | solid dim warm | Dark room |
| `matrix` | digital_rain green | Hacker vibes |
| `waiting` | solid red | Claude needs user input |
| `acknowledged` | breathing green | Claude received input |
| `working` | cycle_spiral | Claude is processing |
| `idle` | breathing blue dim | Nothing happening |

### Custom States

Default states live in `config/default-states.json`. User overrides are stored in `~/.config/saengsation/states.json`.

To add your own states, either:

1. **Edit the defaults** — add entries to `config/default-states.json`
2. **Save from CLI** — set the keyboard how you like it, then `saengsation state save mystate -d "description"`
3. **Edit user config** — create/edit `~/.config/saengsation/states.json`:

```json
{
  "coding": {
    "description": "Purple breathing for late night coding",
    "effect": "breathing",
    "hue": 200,
    "sat": 255,
    "brightness": 100,
    "speed": 1
  }
}
```

User states override defaults with the same name.

## Claude Code Integration

Saengsation includes a Claude Code skill and hooks that change your keyboard lighting based on what Claude is doing.

### How It Looks

| Claude Activity | Keyboard State | Visual |
|----------------|---------------|--------|
| Waiting for your input | `waiting` | Solid red |
| Received your input | `acknowledged` | Green breathing fade |
| Working (tool calls, thinking) | `working` | Rainbow spiral |
| Idle / stopped | `idle` | Slow dim blue pulse |

### Setup Hooks

```bash
make hooks
```

This merges saengsation hooks into `~/.claude/settings.json` with the correct paths for your clone. Restart Claude Code for hooks to take effect.

Requires `jq` if you have existing settings (to safely merge). A template is also available at `skill/claude-settings-example.json` for manual setup.

### Install the Skill

Copy the skill into your project's `.claude/skills/` directory:

```bash
mkdir -p .claude/skills/saengsation
cp skill/saengsation.md .claude/skills/saengsation/SKILL.md
```

Or for global availability across all projects:

```bash
mkdir -p ~/.claude/skills/saengsation
cp skill/saengsation.md ~/.claude/skills/saengsation/SKILL.md
```

### Example Prompts

Once the skill is installed, you can tell Claude:

- *"set my keyboard to focus mode"*
- *"make it red"*
- *"party mode"*
- *"turn off my keyboard lights"*
- *"pulse blue for 10 seconds"*
- *"run police lights"*
- *"save the current keyboard state as 'review'"*

### Customize the Lifecycle

Edit the states in `config/default-states.json` to change what each Claude activity looks like. For example, to make the "working" state a green wave instead of rainbow spiral:

```json
{
  "working": {
    "description": "Claude is working — green wave",
    "effect": "cycle_left_right",
    "hue": 85,
    "sat": 255,
    "brightness": 200,
    "speed": 3
  }
}
```

## Project Structure

```
saengsation/
├── config/
│   └── default-states.json      # Built-in state definitions
├── saengsation/
│   ├── __init__.py
│   ├── __main__.py
│   ├── cli.py                   # CLI entry point
│   ├── keychron.py              # Keychron V7 VIA v10 HID protocol
│   └── states.py                # State loading/saving
├── scripts/
│   ├── check-deps.sh            # Verify deps and permissions
│   ├── setup.sh                 # Full setup
│   ├── install-hidraw.sh        # Build hidapi with hidraw backend
│   ├── install-hooks.sh         # Install Claude Code hooks
│   ├── install.sh               # Install Python deps
│   ├── demo.sh                  # Demo animations
│   └── demo-states.sh           # Demo all states
├── skill/
│   ├── saengsation.md           # Claude Code skill definition
│   ├── hooks.sh                 # Claude Code event hooks
│   └── claude-settings-example.json
├── 99-saengsation.rules         # udev rules
├── Makefile
└── pyproject.toml
```

## Make Targets

```
make help          Show all targets
make setup         Full setup (group, udev, deps, hidraw backend)
make setup-udev    Install udev rules only
make setup-group   Create plugdev group and add user
make install       Install Python dependencies (with hidraw backend)
make hooks         Install Claude Code hooks into ~/.claude/settings.json
make check         Verify dependencies and device access
make demo          Run demo animations
make status        Show keyboard status
make clean         Remove build artifacts
```

## Technical Notes

- Communicates via QMK VIA protocol v10 (0x0A) over raw HID (interface 1, usage page 0xFF60)
- Colors use HSV (hue 0-255, saturation 0-255). Brightness is separate (0-255). Speed is 0-3.
- `--save` persists to keyboard EEPROM. Without it, settings revert on unplug.
- Per-key RGB is not available via the stock VIA protocol (would require custom QMK firmware).
