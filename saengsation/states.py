"""Named lighting states (presets) for the Keychron V7.

Default states are loaded from config/default-states.json in the package.
User overrides are stored in ~/.config/saengsation/states.json.
"""

import json
import os
from pathlib import Path

_PACKAGE_DIR = Path(__file__).parent.parent
DEFAULT_STATES_FILE = _PACKAGE_DIR / "config" / "default-states.json"

CONFIG_DIR = Path(os.environ.get("XDG_CONFIG_HOME", Path.home() / ".config")) / "saengsation"
USER_STATES_FILE = CONFIG_DIR / "states.json"


def _load_default_states() -> dict:
    if DEFAULT_STATES_FILE.exists():
        with open(DEFAULT_STATES_FILE) as f:
            return json.load(f)
    return {}


def _load_user_states() -> dict:
    if USER_STATES_FILE.exists():
        with open(USER_STATES_FILE) as f:
            return json.load(f)
    return {}


def _save_user_states(states: dict):
    CONFIG_DIR.mkdir(parents=True, exist_ok=True)
    with open(USER_STATES_FILE, "w") as f:
        json.dump(states, f, indent=2)


def get_all_states() -> dict:
    merged = _load_default_states()
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
