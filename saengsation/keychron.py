"""Keychron V7 RGB control via QMK VIA raw HID protocol.

VIA protocol v10 (0x0A) uses flat lighting commands with no channel byte.
Format: [command_id, value_id, data...]

Value IDs (confirmed by hardware probing):
  0x80 = brightness (0-255)
  0x81 = effect/mode
  0x82 = speed (0-3)
  0x83 = color (hue, saturation as 2 bytes)
"""

import hid
import time

VID = 0x3434
PID = 0x0370
RAW_HID_INTERFACE = 1
RAW_HID_BUFFER_SIZE = 32

# VIA protocol command IDs
VIA_SET = 0x07
VIA_GET = 0x08
VIA_SAVE = 0x09

# VIA v10 lighting value IDs
VAL_BRIGHTNESS = 0x80
VAL_EFFECT = 0x81
VAL_SPEED = 0x82
VAL_COLOR = 0x83  # 2 bytes: hue, saturation

# Known QMK RGB Matrix effects (common subset, firmware may support fewer)
EFFECTS = {
    "off": 0,
    "solid": 1,
    "breathing": 2,
    "band_spiral_val": 3,
    "cycle_all": 4,
    "cycle_left_right": 5,
    "cycle_up_down": 6,
    "rainbow_moving_chevron": 7,
    "cycle_out_in": 8,
    "cycle_out_in_dual": 9,
    "cycle_pinwheel": 10,
    "cycle_spiral": 11,
    "dual_beacon": 12,
    "rainbow_beacon": 13,
    "jellybean_raindrops": 14,
    "pixel_rain": 15,
    "typing_heatmap": 16,
    "digital_rain": 17,
    "reactive_simple": 18,
    "reactive_multiwide": 19,
    "reactive_multinexus": 20,
    "splash": 21,
    "solid_splash": 22,
}


class KeychronV7:
    def __init__(self):
        self._dev = None

    def open(self):
        self._dev = hid.device()
        path = self._find_raw_hid_path()
        if path is None:
            raise RuntimeError(
                "Keychron V7 not found. Is it plugged in? Are udev rules installed?"
            )
        self._dev.open_path(path)
        self._dev.set_nonblocking(True)

    def close(self):
        if self._dev:
            self._dev.close()
            self._dev = None

    def __enter__(self):
        self.open()
        return self

    def __exit__(self, *args):
        self.close()

    def _find_raw_hid_path(self) -> bytes | None:
        for dev in hid.enumerate(VID, PID):
            if dev["interface_number"] == RAW_HID_INTERFACE:
                return dev["path"]
        return None

    def _send(self, data: list[int]) -> list[int]:
        buf = [0x00] + data + [0x00] * (RAW_HID_BUFFER_SIZE - len(data))
        self._dev.write(buf)
        time.sleep(0.02)
        resp = self._dev.read(RAW_HID_BUFFER_SIZE, timeout_ms=200)
        return resp if resp else []

    def _set_value(self, value_id: int, *data: int):
        self._send([VIA_SET, value_id, *data])

    def _get_value(self, value_id: int, num_bytes: int = 1) -> list[int] | int | None:
        resp = self._send([VIA_GET, value_id])
        if resp and len(resp) >= 2 + num_bytes:
            if num_bytes == 1:
                return resp[2]
            return list(resp[2:2 + num_bytes])
        return None

    def save(self):
        self._send([VIA_SAVE])

    def set_effect(self, effect: int | str):
        if isinstance(effect, str):
            effect = EFFECTS.get(effect.lower(), 1)
        self._set_value(VAL_EFFECT, effect)

    def set_speed(self, speed: int):
        self._set_value(VAL_SPEED, max(0, min(3, speed)))

    def set_brightness(self, val: int):
        self._set_value(VAL_BRIGHTNESS, max(0, min(255, val)))

    def set_color(self, hue: int, sat: int):
        self._set_value(VAL_COLOR, max(0, min(255, hue)), max(0, min(255, sat)))

    def set_hue(self, hue: int):
        color = self._get_value(VAL_COLOR, 2)
        sat = color[1] if color else 255
        self.set_color(hue, sat)

    def set_saturation(self, sat: int):
        color = self._get_value(VAL_COLOR, 2)
        hue = color[0] if color else 0
        self.set_color(hue, sat)

    def set_color_hsv(self, hue: int, sat: int, brightness: int):
        self.set_color(hue, sat)
        self.set_brightness(brightness)

    def get_state(self) -> dict:
        color = self._get_value(VAL_COLOR, 2)
        return {
            "brightness": self._get_value(VAL_BRIGHTNESS),
            "effect": self._get_value(VAL_EFFECT),
            "speed": self._get_value(VAL_SPEED),
            "hue": color[0] if color else None,
            "saturation": color[1] if color else None,
        }

    def cycle_colors(self, duration: float = 10.0, step_delay: float = 0.05):
        self.set_effect("solid")
        self.set_color(0, 255)
        self.set_brightness(255)
        start = time.time()
        hue = 0
        while time.time() - start < duration:
            self.set_color(hue % 256, 255)
            hue += 3
            time.sleep(step_delay)

    def pulse(self, hue: int = 0, cycles: int = 5, step_delay: float = 0.03):
        self.set_effect("solid")
        self.set_color(hue, 255)
        for _ in range(cycles):
            for brightness in range(0, 255, 5):
                self.set_brightness(brightness)
                time.sleep(step_delay)
            for brightness in range(255, 0, -5):
                self.set_brightness(brightness)
                time.sleep(step_delay)

    def flash_alert(self, hue: int = 0, flashes: int = 5, on_time: float = 0.15, off_time: float = 0.1):
        self.set_effect("solid")
        self.set_color(hue, 255)
        for _ in range(flashes):
            self.set_brightness(255)
            time.sleep(on_time)
            self.set_brightness(0)
            time.sleep(off_time)

    def police_lights(self, duration: float = 10.0, flash_time: float = 0.12):
        self.set_effect("solid")
        self.set_brightness(255)
        start = time.time()
        while time.time() - start < duration:
            self.set_color(0, 255)  # red
            time.sleep(flash_time)
            self.set_brightness(0)
            time.sleep(flash_time * 0.5)
            self.set_brightness(255)
            self.set_color(170, 255)  # blue
            time.sleep(flash_time)
            self.set_brightness(0)
            time.sleep(flash_time * 0.5)
            self.set_brightness(255)
