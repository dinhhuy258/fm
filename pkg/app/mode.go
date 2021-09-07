package app

type Action struct {
	help     string
	messages []func(app *App) error
}

type KeyBindings struct {
	onKeys map[string]*Action
}

type Mode struct {
	name        string
	keyBindings *KeyBindings
}
