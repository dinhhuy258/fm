package app

import "github.com/dinhhuy258/fm/pkg/gui"

type Action struct {
	help     string
	messages []func(gui *gui.Gui) error
}

type KeyBindings struct {
	onKeys map[string]*Action
}

type Mode struct {
	name        string
	keyBindings *KeyBindings
}
