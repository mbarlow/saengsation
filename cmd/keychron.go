package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	VID              = 0x3434
	PID              = 0x0370
	RawHIDInterface  = 1
	RawHIDBufferSize = 32

	ViaSet  = 0x07
	ViaGet  = 0x08
	ViaSave = 0x09

	ValBrightness = 0x80
	ValEffect     = 0x81
	ValSpeed      = 0x82
	ValColor      = 0x83 // 2 bytes: hue, saturation
)

// KeychronV7 controls the keyboard RGB LEDs via raw HID over /dev/hidraw.
type KeychronV7 struct {
	file *os.File
}

// findHidrawDevice scans sysfs for the Keychron V7 raw HID interface.
func findHidrawDevice() (string, error) {
	matches, err := filepath.Glob("/sys/class/hidraw/hidraw*")
	if err != nil {
		return "", err
	}

	hidID := fmt.Sprintf("0003:%08X:%08X", VID, PID)

	for _, entry := range matches {
		uevent, err := os.ReadFile(filepath.Join(entry, "device", "uevent"))
		if err != nil {
			continue
		}

		found := false
		for _, line := range strings.Split(string(uevent), "\n") {
			if strings.TrimSpace(line) == "HID_ID="+hidID {
				found = true
				break
			}
		}
		if !found {
			continue
		}

		// Resolve the device symlink to find the USB interface directory.
		devicePath, err := filepath.EvalSymlinks(filepath.Join(entry, "device"))
		if err != nil {
			continue
		}

		// Walk up to find bInterfaceNumber.
		// The device path is like .../usb1/1-3/1-3:1.1/0003:3434:0370.xxxx
		// The USB interface dir (with bInterfaceNumber) is the parent.
		ifDir := filepath.Dir(devicePath)
		ifNum, err := os.ReadFile(filepath.Join(ifDir, "bInterfaceNumber"))
		if err != nil {
			continue
		}

		if strings.TrimSpace(string(ifNum)) == fmt.Sprintf("%02d", RawHIDInterface) {
			devName := filepath.Base(entry)
			return "/dev/" + devName, nil
		}
	}

	return "", fmt.Errorf("Keychron V7 not found. Is it plugged in? Are udev rules installed?")
}

func (kb *KeychronV7) Open() error {
	path, err := findHidrawDevice()
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	kb.file = f
	return nil
}

func (kb *KeychronV7) Close() {
	if kb.file != nil {
		kb.file.Close()
		kb.file = nil
	}
}

func (kb *KeychronV7) send(data []byte) ([]byte, error) {
	buf := make([]byte, RawHIDBufferSize)
	copy(buf, data)

	_, err := kb.file.Write(buf)
	if err != nil {
		return nil, fmt.Errorf("hid write: %w", err)
	}

	time.Sleep(20 * time.Millisecond)

	kb.file.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	resp := make([]byte, RawHIDBufferSize)
	n, err := kb.file.Read(resp)
	if err != nil {
		// Timeout is not fatal — some commands don't respond.
		if os.IsTimeout(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("hid read: %w", err)
	}
	return resp[:n], nil
}

func (kb *KeychronV7) setValue(valueID byte, data ...byte) {
	msg := []byte{ViaSet, valueID}
	msg = append(msg, data...)
	kb.send(msg)
}

func (kb *KeychronV7) getValue(valueID byte, numBytes int) []byte {
	resp, err := kb.send([]byte{ViaGet, valueID})
	if err != nil || resp == nil || len(resp) < 2+numBytes {
		return nil
	}
	return resp[2 : 2+numBytes]
}

func (kb *KeychronV7) Save() {
	kb.send([]byte{ViaSave})
}

func (kb *KeychronV7) SetEffect(effect string) {
	if idx, ok := EffectByName(effect); ok {
		kb.setValue(ValEffect, byte(idx))
	} else {
		// Default to solid if unknown.
		kb.setValue(ValEffect, 1)
	}
}

func (kb *KeychronV7) SetEffectNum(effect int) {
	kb.setValue(ValEffect, byte(effect))
}

func (kb *KeychronV7) SetSpeed(speed int) {
	kb.setValue(ValSpeed, byte(clamp(speed, 0, 3)))
}

func (kb *KeychronV7) SetBrightness(val int) {
	kb.setValue(ValBrightness, byte(clamp(val, 0, 255)))
}

func (kb *KeychronV7) SetColor(hue, sat int) {
	kb.setValue(ValColor, byte(clamp(hue, 0, 255)), byte(clamp(sat, 0, 255)))
}

func (kb *KeychronV7) SetColorHSV(hue, sat, brightness int) {
	kb.SetColor(hue, sat)
	kb.SetBrightness(brightness)
}

func (kb *KeychronV7) GetState() map[string]interface{} {
	state := map[string]interface{}{
		"brightness": nil,
		"effect":     nil,
		"speed":      nil,
		"hue":        nil,
		"saturation": nil,
	}

	if v := kb.getValue(ValBrightness, 1); v != nil {
		state["brightness"] = int(v[0])
	}
	if v := kb.getValue(ValEffect, 1); v != nil {
		state["effect"] = int(v[0])
	}
	if v := kb.getValue(ValSpeed, 1); v != nil {
		state["speed"] = int(v[0])
	}
	if v := kb.getValue(ValColor, 2); v != nil {
		state["hue"] = int(v[0])
		state["saturation"] = int(v[1])
	}

	return state
}

// Animations

func (kb *KeychronV7) CycleColors(duration float64) {
	kb.SetEffect("solid")
	kb.SetColor(0, 255)
	kb.SetBrightness(255)
	start := time.Now()
	hue := 0.0
	for time.Since(start).Seconds() < duration {
		kb.SetColor(int(hue*255)%256, 255)
		hue = hue + 0.01
		if hue >= 1.0 {
			hue -= 1.0
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func (kb *KeychronV7) PoliceLights(duration float64) {
	kb.SetEffect("solid")
	kb.SetBrightness(255)
	start := time.Now()
	for time.Since(start).Seconds() < duration {
		kb.SetColor(0, 255) // red
		time.Sleep(120 * time.Millisecond)
		kb.SetBrightness(0)
		time.Sleep(60 * time.Millisecond)
		kb.SetBrightness(255)
		kb.SetColor(170, 255) // blue
		time.Sleep(120 * time.Millisecond)
		kb.SetBrightness(0)
		time.Sleep(60 * time.Millisecond)
		kb.SetBrightness(255)
	}
}

func (kb *KeychronV7) Pulse(hue int, duration float64) {
	kb.SetEffect("solid")
	kb.SetColor(hue, 255)
	start := time.Now()
	for time.Since(start).Seconds() < duration {
		for b := 0; b < 255; b += 5 {
			if time.Since(start).Seconds() >= duration {
				return
			}
			kb.SetBrightness(b)
			time.Sleep(30 * time.Millisecond)
		}
		for b := 255; b > 0; b -= 5 {
			if time.Since(start).Seconds() >= duration {
				return
			}
			kb.SetBrightness(b)
			time.Sleep(30 * time.Millisecond)
		}
	}
}

func (kb *KeychronV7) FlashAlert(hue int, duration float64) {
	kb.SetEffect("solid")
	kb.SetColor(hue, 255)
	start := time.Now()
	for time.Since(start).Seconds() < duration {
		kb.SetBrightness(255)
		time.Sleep(150 * time.Millisecond)
		kb.SetBrightness(0)
		time.Sleep(100 * time.Millisecond)
	}
}

func clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}
