package tui

import (
	"github.com/charmbracelet/bubbles/key"

	"github.com/dinhhuy258/fm/pkg/actions"
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/pipe"
)

// Model represents the root application state
type Model struct {
	// Core application state
	currentPath string
	termWidth   int
	termHeight  int
	ready       bool

	// Separate models for each component (single source of truth)
	explorerModel    *ExplorerModel
	interactiveModel *InteractiveAreaModel
	helpModel        *HelpModel
	statusBar        StatusBarModel

	// App dependencies
	config *config.Config
	pipe   *pipe.Pipe

	// UI state
	err error

	// Dynamic mode and keybinding system
	modeManager     *ModeManager
	dynamicKeyMap   *DynamicKeyMap
	messageExecutor *actions.ActionHandler
	inputBuffer     string

	// Legacy key bindings (minimal fallback for emergency keys only)
	keys KeyMap
}

// StatusBarModel represents the status bar state
type StatusBarModel struct {
	mode     string
	path     string
	selected int
	total    int
}

// KeyMap defines the key bindings for the application
type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Enter    key.Binding
	Space    key.Binding
	Quit     key.Binding
	Help     key.Binding
	Back     key.Binding
	Home     key.Binding
	End      key.Binding
	PageUp   key.Binding
	PageDown key.Binding
	Tab      key.Binding
	Escape   key.Binding
}

// DefaultKeyMap returns the default key bindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("↓/j", "move down"),
		),
		Left: key.NewBinding(
			key.WithKeys("h", "left"),
			key.WithHelp("←/h", "parent directory"),
		),
		Right: key.NewBinding(
			key.WithKeys("l", "right", "enter"),
			key.WithHelp("→/l/enter", "enter directory"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "open/enter"),
		),
		Space: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "select/toggle"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q/ctrl+c", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Back: key.NewBinding(
			key.WithKeys("backspace", "ctrl+h"),
			key.WithHelp("backspace", "back"),
		),
		Home: key.NewBinding(
			key.WithKeys("home", "g"),
			key.WithHelp("home/g", "go to top"),
		),
		End: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("end/G", "go to bottom"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup", "ctrl+u"),
			key.WithHelp("pgup", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdown", "ctrl+d"),
			key.WithHelp("pgdown", "page down"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "switch focus"),
		),
		Escape: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel/escape"),
		),
	}
}

// ShortHelp returns the short help for the key bindings
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Space, k.Help, k.Quit}
}

// FullHelp returns the full help for the key bindings
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Enter, k.Space, k.Back, k.Tab},
		{k.Home, k.End, k.PageUp, k.PageDown},
		{k.Help, k.Quit, k.Escape},
	}
}

// NewModel creates a new root model
func NewModel(cfg *config.Config, pipe *pipe.Pipe) Model {
	// Initialize separate models for each component
	explorerModel := NewExplorerModel()
	interactiveModel := NewInteractiveAreaModel()
	helpModel := NewHelpModel()

	// Initialize dynamic mode system
	modeManager := NewModeManager(cfg)
	dynamicKeyMap := NewDynamicKeyMap(modeManager)

	// Initialize status bar
	statusBar := StatusBarModel{
		mode:     modeManager.GetCurrentMode(),
		path:     "",
		selected: 0,
		total:    0,
	}

	// Create message executor
	messageExecutor := actions.NewActionHandler(modeManager, pipe)

	return Model{
		currentPath:      "",
		ready:            false,
		explorerModel:    explorerModel,
		interactiveModel: interactiveModel,
		helpModel:        helpModel,
		statusBar:        statusBar,
		config:           cfg,
		pipe:             pipe,
		modeManager:      modeManager,
		dynamicKeyMap:    dynamicKeyMap,
		messageExecutor:  messageExecutor,
		inputBuffer:      "",
		keys:             DefaultKeyMap(), // Emergency fallback only (quit, help, esc)
	}
}

// GetCurrentMode returns the current mode
func (m Model) GetCurrentMode() string {
	return m.modeManager.GetCurrentMode()
}

// SwitchMode switches to a new mode
func (m *Model) SwitchMode(modeName string) error {
	err := m.modeManager.SwitchToMode(modeName)
	if err != nil {
		return err
	}

	// Update status bar
	m.statusBar.mode = m.modeManager.GetCurrentMode()

	// Clear dynamic keymap cache when mode changes
	m.dynamicKeyMap.ClearCache()

	// Handle input state when switching modes
	if modeName == "default" {
		// HideInput now returns a command, but we don't need to handle it here
		// since the mode switch is already happening
		_ = m.interactiveModel.HideInput()
		m.inputBuffer = ""
	}

	return nil
}

// SetInputBuffer sets the input buffer value
func (m *Model) SetInputBuffer(value string) {
	m.inputBuffer = value
}

// GetInputBuffer returns the current input buffer value
func (m Model) GetInputBuffer() string {
	return m.inputBuffer
}

// UpdateInputBufferFromKey updates the input buffer with the last key press
func (m *Model) UpdateInputBufferFromKey(keyStr string) {
	// Handle backspace
	if keyStr == "backspace" {
		if len(m.inputBuffer) > 0 {
			m.inputBuffer = m.inputBuffer[:len(m.inputBuffer)-1]
		}
	} else if len(keyStr) == 1 {
		// For single character keys, append to buffer
		m.inputBuffer += keyStr
	}
}
