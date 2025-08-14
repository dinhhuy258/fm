package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/dinhhuy258/fm/pkg/actions"
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/pipe"
)

// InteractiveAreaMode represents the current mode of the interactive area
type InteractiveAreaMode int

const (
	InteractiveModeNotification InteractiveAreaMode = iota // Show notifications
	InteractiveModeInput                                   // Show input field
	InteractiveModeBuffer                                  // Buffer mode for simple character accumulation
)

// Model represents the fm application state
type Model struct {
	// Core application state
	currentPath string
	termWidth   int
	termHeight  int
	ready       bool

	// Models for fm components
	explorerModel      *ExplorerModel
	notificationModel  *NotificationModel
	inputModel         *InputModel
	helpModel          *HelpModel

	// App dependencies
	config *config.Config
	pipe   *pipe.Pipe

	// UI state
	err                   error
	interactiveMode       InteractiveAreaMode // Track current interactive area mode

	// Dynamic mode and keybinding system
	modeManager     *ModeManager
	keyManager      *KeyManager
	messageExecutor *actions.ActionHandler

	// Legacy key bindings (minimal fallback for emergency keys only)
	keys KeyMap
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
	notificationModel := NewNotificationModel()
	inputModel := NewInputModel()
	helpModel := NewHelpModel()

	// Initialize dynamic mode system
	modeManager := NewModeManager(cfg)
	keyManager := NewKeyManager(modeManager)

	// Create message executor
	messageExecutor := actions.NewActionHandler(modeManager, pipe)

	return Model{
		currentPath:       "",
		ready:             false,
		explorerModel:     explorerModel,
		notificationModel: notificationModel,
		inputModel:        inputModel,
		helpModel:         helpModel,
		config:            cfg,
		pipe:              pipe,
		interactiveMode:   InteractiveModeNotification, // Default to notification mode
		modeManager:       modeManager,
		keyManager:        keyManager,
		messageExecutor:   messageExecutor,
		keys:              DefaultKeyMap(), // Emergency fallback only (quit, help, esc)
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

	// Handle input state when switching modes
	if modeName == "default" {
		// Hide input and show notifications when switching to default mode
		m.interactiveMode = InteractiveModeNotification
		m.inputModel.Hide()
		m.notificationModel.Show()
		m.inputModel.ClearBuffer()
	}

	return nil
}

// SetInputBuffer sets the input buffer value
func (m *Model) SetInputBuffer(value string) {
	m.inputModel.SetBuffer(value)
}

// GetInputBuffer returns the current input buffer value
func (m Model) GetInputBuffer() string {
	return m.inputModel.GetBuffer()
}

// UpdateInputBufferFromKey updates the input buffer with the last key press
func (m *Model) UpdateInputBufferFromKey(keyStr string) {
	m.inputModel.AppendToBuffer(keyStr)
}

// ShowInput switches to input mode and displays the text input field
func (m *Model) ShowInput(initialValue string) {
	m.interactiveMode = InteractiveModeInput
	m.notificationModel.Hide()
	m.inputModel.SetMode(InputModeText)
	m.inputModel.Show(initialValue)
}

// HideInput switches back to notification mode
func (m *Model) HideInput() {
	m.interactiveMode = InteractiveModeNotification
	m.inputModel.Hide()
	m.notificationModel.Show()
}

// SetBufferMode switches to buffer mode for simple character accumulation
func (m *Model) SetBufferMode() {
	m.interactiveMode = InteractiveModeBuffer
	m.notificationModel.Hide()
	m.inputModel.SetMode(InputModeBuffer)
	m.inputModel.Show("")
}

// IsInputMode returns whether currently in input mode
func (m *Model) IsInputMode() bool {
	return m.interactiveMode == InteractiveModeInput
}

// IsBufferMode returns whether currently in buffer mode
func (m *Model) IsBufferMode() bool {
	return m.interactiveMode == InteractiveModeBuffer
}

// GetInputValue returns the current input value
func (m *Model) GetInputValue() string {
	return m.inputModel.GetValue()
}

// ShowNotification displays a notification
func (m *Model) ShowNotification(notificationType NotificationType, message string) tea.Cmd {
	cmd := m.notificationModel.ShowNotification(notificationType, message)

	if m.interactiveMode == InteractiveModeInput || m.interactiveMode == InteractiveModeBuffer {
		// In input/buffer mode, notification is stored but not displayed yet
		return cmd
	}

	// In notification mode, ensure notification is visible
	m.notificationModel.Show()

	return cmd
}

// ShowSuccess displays a success notification (auto-clears in 5 seconds)
func (m *Model) ShowSuccess(message string) tea.Cmd {
	return m.ShowNotification(NotificationSuccess, message)
}

// ShowInfo displays an info notification
func (m *Model) ShowInfo(message string) tea.Cmd {
	return m.ShowNotification(NotificationInfo, message)
}

// ShowWarning displays a warning notification
func (m *Model) ShowWarning(message string) tea.Cmd {
	return m.ShowNotification(NotificationWarning, message)
}

// ShowError displays an error notification
func (m *Model) ShowError(message string) tea.Cmd {
	return m.ShowNotification(NotificationError, message)
}

// GetActiveNotification returns the current active notification
func (m *Model) GetActiveNotification() *Notification {
	return m.notificationModel.GetActiveNotification()
}

// ClearNotification clears the active notification
func (m *Model) ClearNotification() {
	m.notificationModel.ClearNotification()
}

// GetTextInput returns the text input model for direct manipulation
func (m *Model) GetTextInput() *textinput.Model {
	return m.inputModel.GetTextInput()
}

// UpdateTextInput updates the text input model
func (m *Model) UpdateTextInput(textInput textinput.Model) {
	m.inputModel.UpdateTextInput(textInput)
}