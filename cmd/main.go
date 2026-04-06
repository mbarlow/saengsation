package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "status":
		cmdStatus()
	case "effects":
		cmdEffects()
	case "kb":
		cmdKb(os.Args[2:])
	case "state":
		cmdState(os.Args[2:])
	case "animate":
		cmdAnimate(os.Args[2:])
	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage: saengsation <command> [args]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  status     Show keyboard status")
	fmt.Println("  effects    List available effects")
	fmt.Println("  kb         Control keyboard LEDs")
	fmt.Println("  state      Manage named lighting states")
	fmt.Println("  animate    Run animations")
}

func cmdStatus() {
	fmt.Println("=== Keychron V7 ===")
	var kb KeychronV7
	if err := kb.Open(); err != nil {
		fmt.Printf("  unavailable: %v\n", err)
		return
	}
	defer kb.Close()

	state := kb.GetState()
	for _, key := range []string{"brightness", "effect", "speed", "hue", "saturation"} {
		fmt.Printf("  %s: %v\n", key, state[key])
	}
}

func cmdEffects() {
	fmt.Println("Keyboard effects (Keychron V7 / QMK):")
	type entry struct {
		name string
		idx  int
	}
	var entries []entry
	for name, idx := range Effects {
		entries = append(entries, entry{name, idx})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].idx < entries[j].idx })
	for _, e := range entries {
		fmt.Printf("  %3d  %s\n", e.idx, e.name)
	}
}

func cmdKb(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: saengsation kb <effect|color|brightness|speed> <value> [--save]")
		os.Exit(1)
	}

	fs := flag.NewFlagSet("kb", flag.ExitOnError)
	save := fs.Bool("save", false, "Save to EEPROM")
	fs.Parse(args[2:])

	action := args[0]
	value := args[1]

	var kb KeychronV7
	if err := kb.Open(); err != nil {
		fmt.Printf("Keyboard unavailable: %v\n", err)
		os.Exit(1)
	}
	defer kb.Close()

	switch action {
	case "effect":
		if num, err := strconv.Atoi(value); err == nil {
			kb.SetEffectNum(num)
		} else {
			kb.SetEffect(value)
		}
		fmt.Printf("Effect: %s\n", value)

	case "color":
		parts := strings.Split(value, ",")
		nums := make([]int, len(parts))
		for i, p := range parts {
			n, err := strconv.Atoi(strings.TrimSpace(p))
			if err != nil {
				fmt.Println("Expected: hue,sat,brightness or hue,sat (0-255 each)")
				os.Exit(1)
			}
			nums[i] = n
		}
		switch len(nums) {
		case 3:
			kb.SetColorHSV(nums[0], nums[1], nums[2])
			fmt.Printf("HSV: hue=%d sat=%d brightness=%d\n", nums[0], nums[1], nums[2])
		case 2:
			kb.SetColor(nums[0], nums[1])
			fmt.Printf("Color: hue=%d sat=%d\n", nums[0], nums[1])
		default:
			fmt.Println("Expected: hue,sat,brightness or hue,sat (0-255 each)")
			os.Exit(1)
		}

	case "brightness":
		n, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("Expected: 0-255")
			os.Exit(1)
		}
		kb.SetBrightness(n)
		fmt.Printf("Brightness: %s\n", value)

	case "speed":
		n, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("Expected: 0-3")
			os.Exit(1)
		}
		kb.SetSpeed(n)
		fmt.Printf("Speed: %s (0-3)\n", value)

	default:
		fmt.Printf("Unknown action: %s (expected: effect, color, brightness, speed)\n", action)
		os.Exit(1)
	}

	if *save {
		kb.Save()
		fmt.Println("Settings saved to EEPROM.")
	}
}

func cmdState(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: saengsation state <list|set|save|delete|show> [name] [--save] [-d description]")
		os.Exit(1)
	}

	action := args[0]

	switch action {
	case "list":
		states := GetAllStates()
		if len(states) == 0 {
			fmt.Println("No states defined.")
			return
		}
		type entry struct {
			name string
			s    State
		}
		var entries []entry
		for name, s := range states {
			entries = append(entries, entry{name, s})
		}
		sort.Slice(entries, func(i, j int) bool { return entries[i].name < entries[j].name })
		for _, e := range entries {
			fmt.Printf("  %-12s  %-20s  %s\n", e.name, e.s.Effect, e.s.Description)
		}

	case "set":
		fs := flag.NewFlagSet("state-set", flag.ExitOnError)
		save := fs.Bool("save", false, "Save to EEPROM")
		fs.Parse(args[1:])
		if fs.NArg() < 1 {
			fmt.Println("Usage: saengsation state set <name> [--save]")
			os.Exit(1)
		}
		name := fs.Arg(0)
		state, ok := GetState(name)
		if !ok {
			fmt.Printf("Unknown state: %s\n", name)
			fmt.Printf("Available: %s\n", strings.Join(sortedStateNames(), ", "))
			os.Exit(1)
		}
		var kb KeychronV7
		if err := kb.Open(); err != nil {
			fmt.Printf("Keyboard unavailable: %v\n", err)
			os.Exit(1)
		}
		defer kb.Close()
		ApplyState(&kb, state)
		if *save {
			kb.Save()
			fmt.Printf("State '%s' applied and saved to EEPROM.\n", name)
		} else {
			fmt.Printf("State '%s' applied. (%s)\n", name, state.Description)
		}

	case "save":
		fs := flag.NewFlagSet("state-save", flag.ExitOnError)
		desc := fs.String("d", "", "Description")
		fs.StringVar(desc, "description", "", "Description")
		fs.Parse(args[1:])
		if fs.NArg() < 1 {
			fmt.Println("Usage: saengsation state save <name> [-d description]")
			os.Exit(1)
		}
		name := fs.Arg(0)
		var kb KeychronV7
		if err := kb.Open(); err != nil {
			fmt.Printf("Keyboard unavailable: %v\n", err)
			os.Exit(1)
		}
		defer kb.Close()
		state := StateFromKeyboard(&kb, *desc)
		if err := SaveState(name, state); err != nil {
			fmt.Printf("Error saving state: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Saved current keyboard state as '%s'.\n", name)

	case "delete":
		if len(args) < 2 {
			fmt.Println("Usage: saengsation state delete <name>")
			os.Exit(1)
		}
		name := args[1]
		if DeleteState(name) {
			fmt.Printf("Deleted state '%s'.\n", name)
		} else {
			fmt.Printf("'%s' is a built-in state or doesn't exist. Only user states can be deleted.\n", name)
		}

	case "show":
		if len(args) < 2 {
			fmt.Println("Usage: saengsation state show <name>")
			os.Exit(1)
		}
		name := args[1]
		state, ok := GetState(name)
		if !ok {
			fmt.Printf("Unknown state: %s\n", name)
			os.Exit(1)
		}
		fmt.Printf("  description: %s\n", state.Description)
		fmt.Printf("  effect: %s\n", state.Effect)
		fmt.Printf("  hue: %d\n", state.Hue)
		fmt.Printf("  sat: %d\n", state.Sat)
		fmt.Printf("  brightness: %d\n", state.Brightness)
		fmt.Printf("  speed: %d\n", state.Speed)

	default:
		fmt.Printf("Unknown action: %s (expected: list, set, save, delete, show)\n", action)
		os.Exit(1)
	}
}

func cmdAnimate(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: saengsation animate <cycle|police|pulse|flash> [--duration N] [--hue N]")
		os.Exit(1)
	}

	fs := flag.NewFlagSet("animate", flag.ExitOnError)
	duration := fs.Float64("duration", 10.0, "Duration in seconds")
	hue := fs.Int("hue", 0, "Base hue 0-255")
	fs.Parse(args[1:])

	animation := args[0]

	var kb KeychronV7
	if err := kb.Open(); err != nil {
		fmt.Printf("Keyboard unavailable: %v\n", err)
		os.Exit(1)
	}
	defer kb.Close()

	// Handle Ctrl+C gracefully.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	done := make(chan struct{})

	go func() {
		select {
		case <-sig:
			fmt.Println("\nStopped.")
			kb.Close()
			os.Exit(0)
		case <-done:
		}
	}()

	switch animation {
	case "cycle":
		fmt.Printf("Color cycling for %.0fs... (Ctrl+C to stop)\n", *duration)
		kb.CycleColors(*duration)
	case "police":
		fmt.Printf("Police lights for %.0fs... (Ctrl+C to stop)\n", *duration)
		kb.PoliceLights(*duration)
	case "pulse":
		fmt.Printf("Pulsing for %.0fs... (Ctrl+C to stop)\n", *duration)
		kb.Pulse(*hue, *duration)
	case "flash":
		fmt.Printf("Flashing for %.0fs... (Ctrl+C to stop)\n", *duration)
		kb.FlashAlert(*hue, *duration)
	default:
		fmt.Printf("Unknown animation: %s (expected: cycle, police, pulse, flash)\n", animation)
		os.Exit(1)
	}

	close(done)
}

func sortedStateNames() []string {
	all := GetAllStates()
	names := make([]string, 0, len(all))
	for name := range all {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
