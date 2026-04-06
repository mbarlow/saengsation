# saengsation

RGB LED controller for Keychron V7 on Linux, with Claude Code integration.

*แสง (saeng) = light + sensation*

Control your keyboard lighting from the command line, define named states for different activities, and let Claude Code change your keyboard colors based on what it's doing.

## Quickstart

```bash
# Clone
git clone git@github.com:mbarlow/saengsation.git
cd saengsation

# Setup (creates plugdev group, installs udev rules, installs deps)
make setup

# Log out and back in for group permissions, or:
newgrp plugdev

# Verify
make check

# Try it
saengsation state set focus
```

## Install

### Requirements

- Linux (tested on Arch)
- Python 3.12+
- [uv](https://docs.astral.sh/uv/) package manager
- `hidapi` system library (`sudo pacman -S hidapi`)
- Keychron V7 keyboard (USB, QMK firmware)

### Step by Step

```bash
# Install Python deps
uv sync

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

## Usage

### Direct Control

```bash
# Show keyboard status
saengsation status

# Set effect
saengsation kb effect breathing
saengsation kb effect digital_rain
saengsation kb effect 5

# Set color (hue, saturation, brightness — each 0-255)
saengsation kb color 85,255,200       # green
saengsation kb color 0,255,255        # red
saengsation kb color 170,255          # blue (keeps current brightness)

# Brightness and speed
saengsation kb brightness 128
saengsation kb speed 2                # 0-3

# Save to EEPROM (persists across unplugs)
saengsation kb effect breathing --save

# List all effects
saengsation effects
```

### Animations

```bash
saengsation animate cycle                # rainbow color cycle
saengsation animate police               # red/blue flash
saengsation animate pulse --hue 170      # blue pulse
saengsation animate flash --hue 0        # red strobe
saengsation animate cycle --duration 30  # 30 seconds
```

### Named States

States are presets that bundle an effect, color, brightness, and speed into a single name.

```bash
# List all states
saengsation state list

# Apply a state
saengsation state set focus
saengsation state set matrix

# Apply and save to EEPROM
saengsation state set night --save

# View state details
saengsation state show focus

# Save current keyboard settings as a new state
saengsation state save mystate -d "Purple haze for coding"

# Delete a custom state
saengsation state delete mystate
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

Add to your Claude Code settings (`~/.claude/settings.json`). Merge the `hooks` key with any existing settings:

```json
{
  "hooks": {
    "UserPromptSubmit": [
      {
        "matcher": "",
        "hooks": [
          {
            "type": "command",
            "command": "/absolute/path/to/saengsation/skill/hooks.sh acknowledged"
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
            "command": "/absolute/path/to/saengsation/skill/hooks.sh working"
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
            "command": "/absolute/path/to/saengsation/skill/hooks.sh idle"
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
            "command": "/absolute/path/to/saengsation/skill/hooks.sh waiting"
          }
        ]
      }
    ]
  }
}
```

Replace `/absolute/path/to/saengsation` with your actual clone path. A template is in `skill/claude-settings-example.json`.

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
make setup         Full setup (group, udev, deps)
make setup-udev    Install udev rules only
make setup-group   Create plugdev group and add user
make install       Install Python dependencies
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
