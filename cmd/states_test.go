package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaultStates(t *testing.T) {
	states := loadDefaultStates()
	if len(states) == 0 {
		t.Fatal("loadDefaultStates returned empty map")
	}

	expected := []string{"focus", "alert", "chill", "meeting", "music", "night", "matrix", "waiting", "acknowledged", "working", "idle"}
	for _, name := range expected {
		if _, ok := states[name]; !ok {
			t.Errorf("default state %q missing", name)
		}
	}
}

func TestDefaultStatesHaveValidEffects(t *testing.T) {
	states := loadDefaultStates()
	for name, s := range states {
		if _, ok := EffectByName(s.Effect); !ok {
			t.Errorf("state %q has unknown effect %q", name, s.Effect)
		}
	}
}

func TestDefaultStatesJSON(t *testing.T) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(defaultStatesJSON, &raw); err != nil {
		t.Fatalf("default-states.json is invalid JSON: %v", err)
	}
}

func TestIntOrDefault(t *testing.T) {
	tests := []struct {
		val  interface{}
		def  int
		want int
	}{
		{42, 0, 42},
		{nil, 99, 99},
		{"not an int", 7, 7},
		{0, 5, 0},
	}

	for _, tt := range tests {
		got := intOrDefault(tt.val, tt.def)
		if got != tt.want {
			t.Errorf("intOrDefault(%v, %d) = %d, want %d", tt.val, tt.def, got, tt.want)
		}
	}
}

func TestConfigDirXDG(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/tmp/test-xdg")
	got := configDir()
	want := "/tmp/test-xdg/saengsation"
	if got != want {
		t.Errorf("configDir() = %q, want %q", got, want)
	}
}

func TestConfigDirDefault(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "")
	home, _ := os.UserHomeDir()
	got := configDir()
	want := filepath.Join(home, ".config", "saengsation")
	if got != want {
		t.Errorf("configDir() = %q, want %q", got, want)
	}
}

func TestSaveLoadDeleteUserStates(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmp)

	// Initially no user states.
	states := loadUserStates()
	if states != nil {
		t.Fatalf("expected nil user states, got %v", states)
	}

	// Save a state.
	custom := State{
		Description: "test state",
		Effect:      "solid",
		Hue:         100,
		Sat:         200,
		Brightness:  150,
		Speed:       2,
	}
	if err := SaveState("test", custom); err != nil {
		t.Fatalf("SaveState: %v", err)
	}

	// Load it back.
	got, ok := GetState("test")
	if !ok {
		t.Fatal("GetState(\"test\") not found after save")
	}
	if got != custom {
		t.Errorf("GetState(\"test\") = %+v, want %+v", got, custom)
	}

	// It should appear in GetAllStates alongside defaults.
	all := GetAllStates()
	if _, ok := all["test"]; !ok {
		t.Error("custom state missing from GetAllStates")
	}
	if _, ok := all["focus"]; !ok {
		t.Error("default state 'focus' missing from GetAllStates")
	}

	// Delete it.
	if !DeleteState("test") {
		t.Error("DeleteState(\"test\") returned false")
	}
	if _, ok := GetState("test"); ok {
		t.Error("state 'test' still exists after delete")
	}

	// Deleting a default state should return false.
	if DeleteState("focus") {
		t.Error("DeleteState(\"focus\") should return false for built-in state")
	}
}

func TestUserStateOverridesDefault(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmp)

	override := State{
		Description: "my focus",
		Effect:      "solid",
		Hue:         0,
		Sat:         255,
		Brightness:  255,
		Speed:       0,
	}
	SaveState("focus", override)

	got, ok := GetState("focus")
	if !ok {
		t.Fatal("GetState(\"focus\") not found")
	}
	if got.Description != "my focus" {
		t.Errorf("expected user override, got description %q", got.Description)
	}
}
