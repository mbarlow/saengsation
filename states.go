package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed config/default-states.json
var defaultStatesJSON []byte

// State represents a named keyboard lighting preset.
type State struct {
	Description string `json:"description,omitempty"`
	Effect      string `json:"effect"`
	Hue         int    `json:"hue"`
	Sat         int    `json:"sat"`
	Brightness  int    `json:"brightness"`
	Speed       int    `json:"speed"`
}

func configDir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "saengsation")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "saengsation")
}

func userStatesFile() string {
	return filepath.Join(configDir(), "states.json")
}

func loadDefaultStates() map[string]State {
	var states map[string]State
	json.Unmarshal(defaultStatesJSON, &states)
	return states
}

func loadUserStates() map[string]State {
	data, err := os.ReadFile(userStatesFile())
	if err != nil {
		return nil
	}
	var states map[string]State
	json.Unmarshal(data, &states)
	return states
}

func saveUserStates(states map[string]State) error {
	dir := configDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(states, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(userStatesFile(), append(data, '\n'), 0644)
}

// GetAllStates returns default states merged with user overrides.
func GetAllStates() map[string]State {
	all := loadDefaultStates()
	for name, s := range loadUserStates() {
		all[name] = s
	}
	return all
}

// GetState looks up a state by name.
func GetState(name string) (State, bool) {
	all := GetAllStates()
	s, ok := all[name]
	return s, ok
}

// SaveState saves a state to the user config file.
func SaveState(name string, state State) error {
	states := loadUserStates()
	if states == nil {
		states = make(map[string]State)
	}
	states[name] = state
	return saveUserStates(states)
}

// DeleteState removes a state from the user config file.
// Returns false if the state doesn't exist in user config.
func DeleteState(name string) bool {
	states := loadUserStates()
	if states == nil {
		return false
	}
	if _, ok := states[name]; !ok {
		return false
	}
	delete(states, name)
	saveUserStates(states)
	return true
}

// ApplyState applies a state to the keyboard.
func ApplyState(kb *KeychronV7, s State) {
	kb.SetEffect(s.Effect)
	kb.SetColor(s.Hue, s.Sat)
	kb.SetBrightness(s.Brightness)
	kb.SetSpeed(s.Speed)
}

// StateFromKeyboard reads the current keyboard state and returns it as a State.
func StateFromKeyboard(kb *KeychronV7, description string) State {
	current := kb.GetState()
	s := State{
		Description: description,
		Hue:         intOrDefault(current["hue"], 0),
		Sat:         intOrDefault(current["saturation"], 255),
		Brightness:  intOrDefault(current["brightness"], 200),
		Speed:       intOrDefault(current["speed"], 1),
	}
	if eff, ok := current["effect"]; ok && eff != nil {
		if name, found := EffectNames[eff.(int)]; found {
			s.Effect = name
		} else {
			s.Effect = fmt.Sprintf("%d", eff.(int))
		}
	} else {
		s.Effect = "solid"
	}
	return s
}

func intOrDefault(v interface{}, def int) int {
	if v == nil {
		return def
	}
	if i, ok := v.(int); ok {
		return i
	}
	return def
}
