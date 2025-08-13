package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/dinhhuy258/fm/pkg/config"
)

// NotificationType represents different types of status notifications
type NotificationType int8

const (
	NotificationSuccess NotificationType = iota
	NotificationInfo
	NotificationWarning
	NotificationError
)

// StatusNotification represents a single status notification with timestamp
type StatusNotification struct {
	Type      NotificationType
	Message   string
	CreatedAt time.Time
}

// InteractiveAreaMode represents the current mode of the interactive area
type InteractiveAreaMode int

const (
	InteractiveModeNotification InteractiveAreaMode = iota // Show notifications
	InteractiveModeInput                                   // Show input field
)

// InteractiveAreaStyles holds cached styles for notifications
type InteractiveAreaStyles struct {
	successStyle lipgloss.Style
	infoStyle    lipgloss.Style
	warningStyle lipgloss.Style
	errorStyle   lipgloss.Style
}

// InteractiveAreaModel represents the pure state for input and notifications
type InteractiveAreaModel struct {
	// Display dimensions
	width  int
	height int

	// Current mode determines what to display
	currentMode InteractiveAreaMode

	// Input state
	textInput textinput.Model

	// Notification state
	activeNotification *StatusNotification

	// Cached notification styles for performance
	styles *InteractiveAreaStyles
}

// NewInteractiveAreaModel creates a new interactive area model
func NewInteractiveAreaModel() *InteractiveAreaModel {
	// Initialize text input
	ti := textinput.New()
	ti.Prompt = "> "

	model := &InteractiveAreaModel{
		currentMode:        InteractiveModeNotification, // Default to showing notifications
		textInput:          ti,
		activeNotification: nil,
	}

	// Initialize cached styles
	model.initStyles()

	return model
}

// SetSize updates the model dimensions
func (m *InteractiveAreaModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// GetSize returns the current dimensions
func (m *InteractiveAreaModel) GetSize() (int, int) {
	return m.width, m.height
}

// ShowInput switches to input mode and displays the text input field
func (m *InteractiveAreaModel) ShowInput(initialValue string) {
	m.currentMode = InteractiveModeInput
	m.textInput.SetValue(initialValue)
	m.textInput.SetCursor(len(initialValue))
	m.textInput.Focus()
}

// HideInput switches back to notification mode
func (m *InteractiveAreaModel) HideInput() tea.Cmd {
	m.currentMode = InteractiveModeNotification
	m.textInput.Blur()
	m.textInput.SetValue("")

	// If there's a pending notification, start auto-clear timer if needed
	if m.activeNotification != nil && m.activeNotification.Type == NotificationSuccess {
		return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
			return AutoClearMessage{}
		})
	}

	return nil
}

// GetInputValue returns the current input value
func (m *InteractiveAreaModel) GetInputValue() string {
	return m.textInput.Value()
}

// IsInputMode returns whether the area is currently in input mode
func (m *InteractiveAreaModel) IsInputMode() bool {
	return m.currentMode == InteractiveModeInput
}

// ShowNotification displays a notification immediately or stores it for later if in input mode
func (m *InteractiveAreaModel) ShowNotification(notificationType NotificationType, message string) tea.Cmd {
	notification := &StatusNotification{
		Type:      notificationType,
		Message:   message,
		CreatedAt: time.Now(),
	}

	// Always store the notification - it will be displayed based on current mode
	m.activeNotification = notification

	if m.currentMode == InteractiveModeInput {
		// In input mode, notification is stored but not displayed yet
		return nil
	}

	// In notification mode, display immediately and auto-clear if success
	if notificationType == NotificationSuccess {
		return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
			return AutoClearMessage{}
		})
	}

	return nil
}

// ShowSuccess displays a success notification (auto-clears in 5 seconds)
func (m *InteractiveAreaModel) ShowSuccess(message string) tea.Cmd {
	return m.ShowNotification(NotificationSuccess, message)
}

// ShowInfo displays an info notification
func (m *InteractiveAreaModel) ShowInfo(message string) tea.Cmd {
	return m.ShowNotification(NotificationInfo, message)
}

// ShowWarning displays a warning notification
func (m *InteractiveAreaModel) ShowWarning(message string) tea.Cmd {
	return m.ShowNotification(NotificationWarning, message)
}

// ShowError displays an error notification
func (m *InteractiveAreaModel) ShowError(message string) tea.Cmd {
	return m.ShowNotification(NotificationError, message)
}

// GetCurrentMode returns the current interactive mode
func (m *InteractiveAreaModel) GetCurrentMode() InteractiveAreaMode {
	return m.currentMode
}

// GetActiveNotification returns the current active notification
func (m *InteractiveAreaModel) GetActiveNotification() *StatusNotification {
	return m.activeNotification
}

// ClearNotification clears the active notification
func (m *InteractiveAreaModel) ClearNotification() {
	m.activeNotification = nil
}

// GetTextInput returns the text input model for direct manipulation
func (m *InteractiveAreaModel) GetTextInput() *textinput.Model {
	return &m.textInput
}

// UpdateTextInput updates the text input model
func (m *InteractiveAreaModel) UpdateTextInput(textInput textinput.Model) {
	m.textInput = textInput
}

// AutoClearMessage represents a message to auto-clear the notification
type AutoClearMessage struct{}

// InputCompletedMessage indicates that input has been completed
type InputCompletedMessage struct {
	Value string // The final input value
}

// View renders the interactive area view using cached styles
func (m *InteractiveAreaModel) View(cfg *config.Config) string {
	width, height := m.GetSize()
	if width <= 0 || height <= 0 {
		return ""
	}

	switch m.currentMode {
	case InteractiveModeInput:
		return m.renderInput()
	case InteractiveModeNotification:
		return m.renderNotification(cfg)
	default:
		return ""
	}
}

// renderInput renders the text input field
func (m *InteractiveAreaModel) renderInput() string {
	return m.textInput.View()
}

// renderNotification renders the current notification (if any)
func (m *InteractiveAreaModel) renderNotification(cfg *config.Config) string {
	if m.activeNotification == nil {
		return "" // Empty when no notification
	}

	width, _ := m.GetSize()
	var notificationStyle lipgloss.Style
	var prefix, suffix string

	// Get style and prefix/suffix based on notification type
	switch m.activeNotification.Type {
	case NotificationSuccess:
		notificationStyle = m.styles.successStyle
		prefix = cfg.General.LogInfoUI.Prefix
		suffix = cfg.General.LogInfoUI.Suffix
	case NotificationInfo:
		notificationStyle = m.styles.infoStyle
		prefix = cfg.General.LogInfoUI.Prefix
		suffix = cfg.General.LogInfoUI.Suffix
	case NotificationWarning:
		notificationStyle = m.styles.warningStyle
		prefix = cfg.General.LogWarningUI.Prefix
		suffix = cfg.General.LogWarningUI.Suffix
	case NotificationError:
		notificationStyle = m.styles.errorStyle
		prefix = cfg.General.LogErrorUI.Prefix
		suffix = cfg.General.LogErrorUI.Suffix
	}

	// Format message with prefix and suffix
	message := prefix + m.activeNotification.Message + suffix

	// Truncate if needed
	if len(message) > width {
		message = message[:width-3] + "..."
	}

	return notificationStyle.Render(message)
}

// initStyles initializes cached notification styles
func (m *InteractiveAreaModel) initStyles() {
	m.styles = &InteractiveAreaStyles{
		successStyle: fromStyleConfig(config.AppConfig.General.LogInfoUI.Style),
		infoStyle:    fromStyleConfig(config.AppConfig.General.LogInfoUI.Style),
		warningStyle: fromStyleConfig(config.AppConfig.General.LogWarningUI.Style),
		errorStyle:   fromStyleConfig(config.AppConfig.General.LogErrorUI.Style),
	}
}

// InvalidateStyles clears cached styles (call when config changes)
func (m *InteractiveAreaModel) InvalidateStyles() {
	m.styles = nil
	m.initStyles()
}
