package main

import "testing"

func TestEffectByName(t *testing.T) {
	tests := []struct {
		name    string
		wantIdx int
		wantOK  bool
	}{
		{"solid", 1, true},
		{"off", 0, true},
		{"breathing", 2, true},
		{"digital_rain", 17, true},
		{"SOLID", 1, true},
		{"Breathing", 2, true},
		{"nonexistent", 0, false},
	}

	for _, tt := range tests {
		idx, ok := EffectByName(tt.name)
		if ok != tt.wantOK || idx != tt.wantIdx {
			t.Errorf("EffectByName(%q) = (%d, %v), want (%d, %v)", tt.name, idx, ok, tt.wantIdx, tt.wantOK)
		}
	}
}

func TestEffectNamesReverseMap(t *testing.T) {
	// Every entry in Effects should have a reverse mapping.
	for name, idx := range Effects {
		got, ok := EffectNames[idx]
		if !ok {
			t.Errorf("EffectNames missing index %d (for %q)", idx, name)
		}
		if got != name {
			// Collisions are possible if two names map to the same index,
			// but in the current table there should be none.
			t.Errorf("EffectNames[%d] = %q, expected %q", idx, got, name)
		}
	}
}

func TestEffectsCount(t *testing.T) {
	if len(Effects) != len(EffectNames) {
		t.Errorf("Effects has %d entries but EffectNames has %d", len(Effects), len(EffectNames))
	}
}
