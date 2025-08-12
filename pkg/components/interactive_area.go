package components

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/dinhhuy258/fm/pkg/config"
)

const inputPrompt = "> "

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

// AutoClearMessage represents a message to auto-clear the notification
type AutoClearMessage struct{}

// InputCompletedMessage indicates that input has been completed
type InputCompletedMessage struct {
	Value string // The final input value
}

// InteractiveAreaMode represents the current mode of the interactive area
type InteractiveAreaMode int

const (
	InteractiveModeNotification InteractiveAreaMode = iota // Show notifications
	InteractiveModeInput                                   // Show input field
)

// InteractiveArea combines input and notification functionality in one space
type InteractiveArea struct {
	width  int
	height int

	// Current mode determines what to display
	currentMode InteractiveAreaMode

	// Input components
	textInput textinput.Model

	// Notification components
	activeNotification *StatusNotification
	successStyle       lipgloss.Style
	infoStyle          lipgloss.Style
	warningStyle       lipgloss.Style
	errorStyle         lipgloss.Style

	// Configuration
	config *config.Config
}

// NewInteractiveArea creates a new interactive area that handles both input and notifications
func NewInteractiveArea() *InteractiveArea {
	cfg := config.AppConfig

	// Initialize text input
	ti := textinput.New()
	ti.Prompt = inputPrompt

	ia := &InteractiveArea{
		currentMode:        InteractiveModeNotification, // Default to showing notifications
		textInput:          ti,
		activeNotification: nil,
		config:             cfg,
	}

	ia.initStyles()

	return ia
}

// initStyles initializes lipgloss styles from config
func (ia *InteractiveArea) initStyles() {
	ia.infoStyle = FromStyleConfig(ia.config.General.LogInfoUI.Style)
	ia.successStyle = FromStyleConfig(ia.config.General.LogInfoUI.Style)
	ia.warningStyle = FromStyleConfig(ia.config.General.LogWarningUI.Style)
	ia.errorStyle = FromStyleConfig(ia.config.General.LogErrorUI.Style)
}

// SetSize updates the interactive area size
func (ia *InteractiveArea) SetSize(width, height int) {
	ia.width = width
	ia.height = height
}

// ShowInput switches to input mode and displays the text input field
func (ia *InteractiveArea) ShowInput(initialValue string) {
	ia.currentMode = InteractiveModeInput
	ia.textInput.SetValue(initialValue)
	ia.textInput.SetCursor(len(initialValue))
	ia.textInput.Focus()
}

// HideInput switches back to notification mode and shows any pending notification
func (ia *InteractiveArea) HideInput() tea.Cmd {
	ia.currentMode = InteractiveModeNotification
	ia.textInput.Blur()
	ia.textInput.SetValue("")

	// If there's a pending notification, start auto-clear timer if needed
	if ia.activeNotification != nil && ia.activeNotification.Type == NotificationSuccess {
		return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
			return AutoClearMessage{}
		})
	}

	return nil
}

// GetInputValue returns the current input value
func (ia *InteractiveArea) GetInputValue() string {
	return ia.textInput.Value()
}

// IsInputMode returns whether the area is currently in input mode
func (ia *InteractiveArea) IsInputMode() bool {
	return ia.currentMode == InteractiveModeInput
}

// ShowNotification displays a notification immediately or stores it for later if in input mode
func (ia *InteractiveArea) ShowNotification(notificationType NotificationType, message string) tea.Cmd {
	notification := &StatusNotification{
		Type:      notificationType,
		Message:   message,
		CreatedAt: time.Now(),
	}

	// Always store the notification - it will be displayed based on current mode
	ia.activeNotification = notification

	if ia.currentMode == InteractiveModeInput {
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
func (ia *InteractiveArea) ShowSuccess(message string) tea.Cmd {
	return ia.ShowNotification(NotificationSuccess, message)
}

// ShowInfo displays an info notification
func (ia *InteractiveArea) ShowInfo(message string) tea.Cmd {
	return ia.ShowNotification(NotificationInfo, message)
}

// ShowWarning displays a warning notification
func (ia *InteractiveArea) ShowWarning(message string) tea.Cmd {
	return ia.ShowNotification(NotificationWarning, message)
}

// ShowError displays an error notification
func (ia *InteractiveArea) ShowError(message string) tea.Cmd {
	return ia.ShowNotification(NotificationError, message)
}

// Update handles input messages
func (ia *InteractiveArea) Update(msg tea.Msg) (*InteractiveArea, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case AutoClearMessage:
		// Auto-clear notification only if in notification mode
		if ia.currentMode == InteractiveModeNotification {
			ia.activeNotification = nil
		}

		return ia, nil
	case tea.KeyMsg:
		if ia.currentMode == InteractiveModeInput {
			// Handle input mode keys
			switch msg.String() {
			case "enter":
				// Input completed - return to notification mode and send completion message
				inputValue := ia.textInput.Value()
				cmd = ia.HideInput()

				// Combine the HideInput command with the completion message
				completionCmd := func() tea.Msg {
					return InputCompletedMessage{Value: inputValue}
				}

				if cmd != nil {
					return ia, tea.Batch(cmd, completionCmd)
				}

				return ia, completionCmd
			case "esc":
				// Cancel input - return to notification mode
				cmd = ia.HideInput()

				return ia, cmd
			default:
				// Pass other keys to text input
				ia.textInput, cmd = ia.textInput.Update(msg)

				return ia, cmd
			}
		}
	}

	return ia, nil
}

// View renders the interactive area based on current mode
func (ia *InteractiveArea) View() string {
	if ia.width <= 0 || ia.height <= 0 {
		return ""
	}

	switch ia.currentMode {
	case InteractiveModeInput:
		return ia.renderInput()
	case InteractiveModeNotification:
		return ia.renderNotification()
	default:
		return ""
	}
}

// renderInput renders the text input field
func (ia *InteractiveArea) renderInput() string {
	return ia.textInput.View()
}

// renderNotification renders the current notification (if any)
func (ia *InteractiveArea) renderNotification() string {
	if ia.activeNotification == nil {
		return "" // Empty when no notification
	}

	notification := ia.activeNotification
	var style lipgloss.Style
	var prefix, suffix string

	// Get style and prefix/suffix based on notification type
	switch notification.Type {
	case NotificationSuccess:
		style = ia.successStyle
		prefix = ia.config.General.LogInfoUI.Prefix
		suffix = ia.config.General.LogInfoUI.Suffix
	case NotificationInfo:
		style = ia.infoStyle
		prefix = ia.config.General.LogInfoUI.Prefix
		suffix = ia.config.General.LogInfoUI.Suffix
	case NotificationWarning:
		style = ia.warningStyle
		prefix = ia.config.General.LogWarningUI.Prefix
		suffix = ia.config.General.LogWarningUI.Suffix
	case NotificationError:
		style = ia.errorStyle
		prefix = ia.config.General.LogErrorUI.Prefix
		suffix = ia.config.General.LogErrorUI.Suffix
	}

	// Format message with prefix and suffix
	message := prefix + notification.Message + suffix

	// Truncate if needed
	if len(message) > ia.width {
		message = message[:ia.width-3] + "..."
	}

	return style.Render(message)
}
