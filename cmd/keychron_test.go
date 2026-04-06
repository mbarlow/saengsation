package main

import "testing"

func TestClamp(t *testing.T) {
	tests := []struct {
		val, min, max, want int
	}{
		{128, 0, 255, 128},
		{-5, 0, 255, 0},
		{300, 0, 255, 255},
		{0, 0, 255, 0},
		{255, 0, 255, 255},
		{2, 0, 3, 2},
		{5, 0, 3, 3},
		{-1, 0, 3, 0},
	}

	for _, tt := range tests {
		got := clamp(tt.val, tt.min, tt.max)
		if got != tt.want {
			t.Errorf("clamp(%d, %d, %d) = %d, want %d", tt.val, tt.min, tt.max, got, tt.want)
		}
	}
}

func TestConstants(t *testing.T) {
	if VID != 0x3434 {
		t.Errorf("VID = 0x%04X, want 0x3434", VID)
	}
	if PID != 0x0370 {
		t.Errorf("PID = 0x%04X, want 0x0370", PID)
	}
	if RawHIDBufferSize != 32 {
		t.Errorf("RawHIDBufferSize = %d, want 32", RawHIDBufferSize)
	}
}
