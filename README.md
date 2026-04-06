# saengsation

RGB LED controller for Keychron V7 on Linux.

*แสง (saeng) = light + sensation*

## Supported Device

| Device | Protocol | Features |
|---|---|---|
| Keychron V7 | QMK VIA v10 raw HID | Effects, color (hue/sat), brightness, speed, save to EEPROM |

## Quickstart

```bash
# 1. Clone and enter the project
cd saengsation

# 2. Run full setup (creates plugdev group, installs udev rules, installs deps)
make setup

# 3. Log out and back in, or:
newgrp plugdev

# 4. Verify everything is working
make check

# 5. Try it out
make demo
```

## Usage

All commands below assume an activated venv (`source .venv/bin/activate`).
Alternatively, prefix each command with `uv run` (e.g. `uv run saengsation status`).

```bash
# Show keyboard status
saengsation status

# List all available effects
saengsation effects
```

### Setting Effects

```bash
saengsation kb effect breathing
saengsation kb effect cycle_spiral
saengsation kb effect digital_rain

# By number
saengsation kb effect 5
```

### Setting Color

```bash
# hue,saturation,brightness (each 0-255)
saengsation kb color 85,255,200       # green, full sat, bright
saengsation kb color 0,255,255        # red, full sat, max bright
saengsation kb color 170,255,255      # blue, full sat, max bright

# hue,saturation only (keeps current brightness)
saengsation kb color 85,255           # green, full sat
```

### Brightness and Speed

```bash
saengsation kb brightness 128
saengsation kb speed 2                # 0-3 (slow to fast)
```

### Save to EEPROM

```bash
# Persist current settings across unplugs
saengsation kb effect breathing --save
```

### Animations

```bash
# Rainbow color cycle
saengsation animate cycle

# Police lights (red/blue flash)
saengsation animate police

# Pulse a specific color (hue 0-255)
saengsation animate pulse --hue 170   # blue pulse
saengsation animate pulse --hue 0     # red pulse

# Strobe flash
saengsation animate flash --hue 85    # green strobe

# Set duration
saengsation animate cycle --duration 30
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

## Manual Setup

```bash
uv sync

sudo groupadd plugdev
sudo usermod -aG plugdev $USER

sudo cp 99-saengsation.rules /etc/udev/rules.d/
sudo udevadm control --reload-rules
sudo udevadm trigger

# Log out/in, then:
uv run saengsation status
```

## Requirements

- Linux (tested on Arch)
- Python 3.12+
- [uv](https://docs.astral.sh/uv/) package manager
- `hidapi` system library (`sudo pacman -S hidapi`)

## Notes

- The effect list is the common QMK RGB matrix set. Your V7 firmware may support a subset — unsupported effect numbers are silently ignored.
- Colors use HSV (hue 0-255, saturation 0-255). Brightness is set separately (0-255). Speed is 0-3.
- Use `--save` to persist settings to EEPROM. Without it, settings revert on unplug.
- Per-key RGB control is not available through the stock VIA protocol. It would require custom QMK firmware.
