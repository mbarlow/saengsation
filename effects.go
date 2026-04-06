package main

import "strings"

// QMK RGB Matrix effects supported by Keychron V7.
var Effects = map[string]int{
	"off":                     0,
	"solid":                   1,
	"breathing":               2,
	"band_spiral_val":         3,
	"cycle_all":               4,
	"cycle_left_right":        5,
	"cycle_up_down":           6,
	"rainbow_moving_chevron":  7,
	"cycle_out_in":            8,
	"cycle_out_in_dual":       9,
	"cycle_pinwheel":          10,
	"cycle_spiral":            11,
	"dual_beacon":             12,
	"rainbow_beacon":          13,
	"jellybean_raindrops":     14,
	"pixel_rain":              15,
	"typing_heatmap":          16,
	"digital_rain":            17,
	"reactive_simple":         18,
	"reactive_multiwide":      19,
	"reactive_multinexus":     20,
	"splash":                  21,
	"solid_splash":            22,
}

// Reverse map: index -> name.
var EffectNames map[int]string

func init() {
	EffectNames = make(map[int]string, len(Effects))
	for name, idx := range Effects {
		EffectNames[idx] = name
	}
}

// EffectByName looks up an effect by name (case-insensitive).
func EffectByName(name string) (int, bool) {
	idx, ok := Effects[strings.ToLower(name)]
	return idx, ok
}
