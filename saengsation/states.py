"""Named lighting states (presets) for the Keychron V7."""

import json
import os
from pathlib import Path

DEFAULT_STATES = {
    "focus": {
        "description": "Calm blue breathing for deep work",
        "effect": "breathing",
        "hue": 170,
        "sat": 255,
        "brightness": 120,
        "speed": 1,
    },
    "alert": {
        "description": "Red flash — something needs attention",
        "effect": "solid",
        "hue": 0,
        "sat": 255,
        "brightness": 255,
        "speed": 0,
    },
    "chill": {
        "description": "Slow rainbow cycle",
        "effect": "cycle_all",
        "hue": 0,
        "sat": 255,
        "brightness": 150,
        "speed": 1,
    },
    "meeting": {
        "description": "Lights off — no distractions",
        "effect": "off",
        "hue": 0,
        "sat": 0,
        "brightness": 0,
        "speed": 0,
    },
    "music": {
        "description": "Rainbow beacon party mode",
        "effect": "rainbow_beacon",
        "hue": 0,
        "sat": 255,
        "brightness": 200,
        "speed": 3,
    },
    "night": {
        "description": "Dim warm white for dark rooms",
        "effect": "solid",
        "hue": 20,
        "sat": 180,
        "brightness": 40,
        "speed": 0,
    },
    "matrix": {
        "description": "Green digital rain",
        "effect": "digital_rain",
        "hue": 85,
        "sat": 255,
        "brightness": 200,
        "speed": 2,
    },
}

CONFIG_DIR = Path(os.environ.get("XDG_CONFIG_HOME", Path.home() / ".config")) / "saengsation"
STATES_FILE = CONFIG_DIR / "states.json"


def _load_user_states() -> dict:
    if STATES_FILE.exists():
        with open(STATES_FILE) as f:
            return json.load(f)
    return {}


def _save_user_states(states: dict):
    CONFIG_DIR.mkdir(parents=True, exist_ok=True)
    with open(STATES_FILE, "w") as f:
        json.dump(states, f, indent=2)


def get_all_states() -> dict:
    merged = dict(DEFAULT_STATES)
    merged.update(_load_user_states())
    return merged


def get_state(name: str) -> dict | None:
    return get_all_states().get(name)


def save_state(name: str, state: dict):
    user_states = _load_user_states()
    user_states[name] = state
    _save_user_states(user_states)


def delete_state(name: str) -> bool:
    user_states = _load_user_states()
    if name in user_states:
        del user_states[name]
        _save_user_states(user_states)
        return True
    return False
