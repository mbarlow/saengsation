"""Saengsation CLI — RGB LED control for Keychron V7."""

import argparse
import sys
import time

from saengsation.keychron import KeychronV7, EFFECTS as KB_EFFECTS
from saengsation.states import get_all_states, get_state, save_state, delete_state


def _apply_state(kb: KeychronV7, state: dict):
    kb.set_effect(state.get("effect", "solid"))
    if "hue" in state and "sat" in state:
        kb.set_color(state["hue"], state["sat"])
    if "brightness" in state:
        kb.set_brightness(state["brightness"])
    if "speed" in state:
        kb.set_speed(state["speed"])


def cmd_status(args):
    print("=== Keychron V7 ===")
    try:
        with KeychronV7() as kb:
            state = kb.get_state()
            for k, v in state.items():
                print(f"  {k}: {v}")
    except Exception as e:
        print(f"  unavailable: {e}")


def cmd_kb(args):
    with KeychronV7() as kb:
        if args.action == "effect":
            name = args.value
            if name.isdigit():
                kb.set_effect(int(name))
            else:
                kb.set_effect(name)
            print(f"Effect: {name}")

        elif args.action == "color":
            parts = [int(x) for x in args.value.split(",")]
            if len(parts) == 3:
                kb.set_color_hsv(*parts)
                print(f"HSV: hue={parts[0]} sat={parts[1]} brightness={parts[2]}")
            elif len(parts) == 2:
                kb.set_color(*parts)
                print(f"Color: hue={parts[0]} sat={parts[1]}")
            else:
                print("Expected: hue,sat,brightness or hue,sat (0-255 each)")
                sys.exit(1)

        elif args.action == "brightness":
            kb.set_brightness(int(args.value))
            print(f"Brightness: {args.value}")

        elif args.action == "speed":
            kb.set_speed(int(args.value))
            print(f"Speed: {args.value} (0-3)")

        if args.save:
            kb.save()
            print("Settings saved to EEPROM.")


def cmd_state(args):
    if args.action == "list":
        states = get_all_states()
        if not states:
            print("No states defined.")
            return
        for name, s in sorted(states.items()):
            desc = s.get("description", "")
            effect = s.get("effect", "?")
            print(f"  {name:<12s}  {effect:<20s}  {desc}")

    elif args.action == "set":
        state = get_state(args.name)
        if state is None:
            print(f"Unknown state: {args.name}")
            print(f"Available: {', '.join(sorted(get_all_states()))}")
            sys.exit(1)
        with KeychronV7() as kb:
            _apply_state(kb, state)
            if args.save:
                kb.save()
                print(f"State '{args.name}' applied and saved to EEPROM.")
            else:
                print(f"State '{args.name}' applied. ({state.get('description', '')})")

    elif args.action == "save":
        with KeychronV7() as kb:
            current = kb.get_state()
        new_state = {
            "effect": current.get("effect", 1),
            "hue": current.get("hue", 0),
            "sat": current.get("saturation", 255),
            "brightness": current.get("brightness", 200),
            "speed": current.get("speed", 1),
        }
        if args.description:
            new_state["description"] = args.description
        save_state(args.name, new_state)
        print(f"Saved current keyboard state as '{args.name}'.")

    elif args.action == "delete":
        if delete_state(args.name):
            print(f"Deleted state '{args.name}'.")
        else:
            print(f"'{args.name}' is a built-in state or doesn't exist. Only user states can be deleted.")

    elif args.action == "show":
        state = get_state(args.name)
        if state is None:
            print(f"Unknown state: {args.name}")
            sys.exit(1)
        for k, v in state.items():
            print(f"  {k}: {v}")


def cmd_animate(args):
    try:
        kb = KeychronV7()
        kb.open()
    except Exception as e:
        print(f"Keyboard unavailable: {e}")
        sys.exit(1)

    try:
        duration = args.duration
        anim = args.animation

        if anim == "cycle":
            print(f"Color cycling for {duration}s... (Ctrl+C to stop)")
            kb.set_effect("solid")
            kb.set_color(0, 255)
            kb.set_brightness(255)
            start = time.time()
            hue = 0.0
            while time.time() - start < duration:
                kb.set_color(int(hue * 255) % 256, 255)
                hue = (hue + 0.01) % 1.0
                time.sleep(0.05)

        elif anim == "police":
            print(f"Police lights for {duration}s... (Ctrl+C to stop)")
            kb.set_effect("solid")
            kb.set_brightness(255)
            start = time.time()
            while time.time() - start < duration:
                kb.set_color(0, 255)  # red
                time.sleep(0.12)
                kb.set_brightness(0)
                time.sleep(0.06)
                kb.set_brightness(255)
                kb.set_color(170, 255)  # blue
                time.sleep(0.12)
                kb.set_brightness(0)
                time.sleep(0.06)
                kb.set_brightness(255)

        elif anim == "pulse":
            print(f"Pulsing for {duration}s... (Ctrl+C to stop)")
            kb.set_effect("solid")
            kb.set_color(args.hue, 255)
            start = time.time()
            while time.time() - start < duration:
                for brightness in range(0, 255, 5):
                    if time.time() - start >= duration:
                        break
                    kb.set_brightness(brightness)
                    time.sleep(0.03)
                for brightness in range(255, 0, -5):
                    if time.time() - start >= duration:
                        break
                    kb.set_brightness(brightness)
                    time.sleep(0.03)

        elif anim == "flash":
            print(f"Flashing for {duration}s... (Ctrl+C to stop)")
            kb.set_effect("solid")
            kb.set_color(args.hue, 255)
            start = time.time()
            while time.time() - start < duration:
                kb.set_brightness(255)
                time.sleep(0.15)
                kb.set_brightness(0)
                time.sleep(0.1)

    except KeyboardInterrupt:
        print("\nStopped.")
    finally:
        kb.close()


def cmd_effects(args):
    print("Keyboard effects (Keychron V7 / QMK):")
    for name, idx in sorted(KB_EFFECTS.items(), key=lambda x: x[1]):
        print(f"  {idx:3d}  {name}")


def main():
    parser = argparse.ArgumentParser(
        prog="saengsation",
        description="RGB LED controller for Keychron V7",
    )
    sub = parser.add_subparsers(dest="command")

    sub.add_parser("status", help="Show keyboard status")
    sub.add_parser("effects", help="List available effects")

    kb_p = sub.add_parser("kb", help="Control keyboard LEDs")
    kb_p.add_argument("action", choices=["effect", "color", "brightness", "speed"])
    kb_p.add_argument("value", help="Effect name/number, hue,sat[,brightness], or 0-255")
    kb_p.add_argument("--save", action="store_true", help="Save to EEPROM")

    state_p = sub.add_parser("state", help="Manage named lighting states")
    state_p.add_argument("action", choices=["list", "set", "save", "delete", "show"])
    state_p.add_argument("name", nargs="?", help="State name")
    state_p.add_argument("--save", action="store_true", dest="save", help="Also save to EEPROM")
    state_p.add_argument("--description", "-d", help="Description when saving a state")

    anim_p = sub.add_parser("animate", help="Run animations")
    anim_p.add_argument("animation", choices=["cycle", "police", "pulse", "flash"])
    anim_p.add_argument("--duration", type=float, default=10.0, help="Duration in seconds")
    anim_p.add_argument("--hue", type=int, default=0, help="Base hue 0-255 (for pulse/flash)")

    args = parser.parse_args()

    if args.command is None:
        parser.print_help()
        sys.exit(0)

    commands = {
        "status": cmd_status,
        "effects": cmd_effects,
        "kb": cmd_kb,
        "state": cmd_state,
        "animate": cmd_animate,
    }
    commands[args.command](args)


if __name__ == "__main__":
    main()
