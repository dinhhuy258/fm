package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/config/lua"
	"github.com/dinhhuy258/fm/pkg/pipe"
	"github.com/dinhhuy258/fm/pkg/tui"
)

var version = "unversioned"

func main() {
	showVersion := flag.Bool("version", false, "Print the current version")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	// Initialize Lua configuration
	luaEngine := lua.NewLua()
	defer luaEngine.Close()

	// Load the config
	if err := config.LoadConfig(luaEngine); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize pipe for external commands
	pipe, err := pipe.NewPipe()
	if err != nil {
		log.Fatalf("failed to create pipe: %v", err)
	}
	defer pipe.StopWatcher()

	// Create the Bubble Tea model
	model := tui.NewModel(pipe)

	// Create the Bubble Tea program
	program := tea.NewProgram(
		model,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	// Start the pipe watcher (for external commands)
	pipe.StartWatcher(func(message string) {
		// Send external messages to the TUI
		program.Send(tui.PipeMessage{Command: message})
	})

	// Run the program
	if _, err := program.Run(); err != nil {
		log.Fatalf("Error running Bubble Tea program: %v", err)
	}
}