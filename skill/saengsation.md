---
name: saengsation
description: Control Keychron V7 keyboard RGB LEDs — set lighting states, effects, colors, and animations
---

You have access to the `saengsation` CLI tool for controlling the Keychron V7 keyboard RGB LEDs.

## Available Commands

```bash
# Apply a named state
saengsation state set <name>

# List all states
saengsation state list

# Set effect directly
saengsation kb effect <name|number>

# Set color (HSV)
saengsation kb color <hue>,<sat>[,<brightness>]

# Set brightness / speed
saengsation kb brightness <0-255>
saengsation kb speed <0-3>

# Run an animation
saengsation animate <cycle|police|pulse|flash> [--duration N] [--hue N]
```

## Built-in States

| State | Effect | Description |
|-------|--------|-------------|
| focus | breathing blue | Deep work mode |
| alert | solid red | Something needs attention |
| chill | cycle_all | Slow rainbow cycle |
| meeting | off | Lights off |
| music | rainbow_beacon | Party mode |
| night | solid dim warm | Dark room |
| matrix | digital_rain green | Hacker vibes |
| waiting | solid red | Claude needs user input |
| acknowledged | breathing green | Claude received input |
| working | cycle_spiral | Claude is processing |
| idle | breathing blue dim | Nothing happening |

## When to Use

- When the user asks to change their keyboard lighting
- When the user references a mood, state, or activity that maps to a preset
- When the user asks what lighting options are available

## Examples

User: "set my keyboard to focus mode"
→ `saengsation state set focus`

User: "make it red"
→ `saengsation kb color 0,255,255`

User: "party mode"
→ `saengsation state set music`

User: "turn off my keyboard lights"
→ `saengsation state set meeting`

User: "run police lights for 5 seconds"
→ `saengsation animate police --duration 5`
