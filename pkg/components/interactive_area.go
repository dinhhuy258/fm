package components

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/dinhhuy258/fm/pkg/config"
)

// InputCompletedMessage indicates that input has been completed
// (This is defined here to avoid circular imports with actions package)
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
	textInput   textinput.Model
	inputPrompt string

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
	ti.Prompt = "> "
	ti.CharLimit = 256

	ia := &InteractiveArea{
		currentMode:        InteractiveModeNotification, // Default to showing notifications
		textInput:          ti,
		inputPrompt:        "> ",
		activeNotification: nil,
		config:             cfg,
	}

	ia.initStyles()

	return ia
}

// initStyles initializes lipgloss styles from config
func (ia *InteractiveArea) initStyles() {
	// Default styles with good visibility
	ia.successStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")).
		Bold(true)

	ia.infoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(false)

	ia.warningStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFF00")).
		Bold(true)

	ia.errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000")).
		Bold(true)

	// Override with config if available (reusing existing log UI configs)
	if ia.config != nil && ia.config.General.LogInfoUI.Style != nil {
		ia.infoStyle = convertConfigStyleToLipgloss(ia.config.General.LogInfoUI.Style)
		ia.successStyle = convertConfigStyleToLipgloss(
			ia.config.General.LogInfoUI.Style,
		).Foreground(lipgloss.Color("#00FF00"))
	}
	if ia.config != nil && ia.config.General.LogWarningUI.Style != nil {
		ia.warningStyle = convertConfigStyleToLipgloss(ia.config.General.LogWarningUI.Style)
	}
	if ia.config != nil && ia.config.General.LogErrorUI.Style != nil {
		ia.errorStyle = convertConfigStyleToLipgloss(ia.config.General.LogErrorUI.Style)
	}
}

// SetSize updates the interactive area size
func (ia *InteractiveArea) SetSize(width, height int) {
	ia.width = width
	ia.height = height
}

// --- Input Methods ---

// ShowInput switches to input mode and displays the text input field
func (ia *InteractiveArea) ShowInput(prompt, initialValue string) {
	ia.currentMode = InteractiveModeInput
	if prompt != "" {
		ia.textInput.Prompt = prompt
	}
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

// --- Notification Methods ---

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

// ClearNotification clears the current notification
func (ia *InteractiveArea) ClearNotification() {
	ia.activeNotification = nil
}

// HasActiveNotification returns whether there's an active notification
func (ia *InteractiveArea) HasActiveNotification() bool {
	return ia.currentMode == InteractiveModeNotification && ia.activeNotification != nil
}

// HasPendingNotification returns whether there's a notification waiting to be shown
func (ia *InteractiveArea) HasPendingNotification() bool {
	return ia.currentMode == InteractiveModeInput && ia.activeNotification != nil
}

// --- Update and View ---

// Update handles input messages
func (ia *InteractiveArea) Update(msg tea.Msg) (*InteractiveArea, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case AutoClearMessage:
		// Auto-clear notification only if in notification mode
		if ia.currentMode == InteractiveModeNotification {
			ia.ClearNotification()
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
		} else {
			// Handle notification mode keys
			if msg.String() == "esc" && ia.HasActiveNotification() {
				ia.ClearNotification()

				return ia, nil
			}
		}
	default:
		// Pass other messages to text input if in input mode
		if ia.currentMode == InteractiveModeInput {
			ia.textInput, cmd = ia.textInput.Update(msg)

			return ia, cmd
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
		if ia.config != nil {
			prefix = ia.config.General.LogInfoUI.Prefix
			suffix = ia.config.General.LogInfoUI.Suffix
		}
		if prefix == "" {
			prefix = "✓ "
		}
	case NotificationInfo:
		style = ia.infoStyle
		if ia.config != nil {
			prefix = ia.config.General.LogInfoUI.Prefix
			suffix = ia.config.General.LogInfoUI.Suffix
		}
		if prefix == "" {
			prefix = "ℹ "
		}
	case NotificationWarning:
		style = ia.warningStyle
		if ia.config != nil {
			prefix = ia.config.General.LogWarningUI.Prefix
			suffix = ia.config.General.LogWarningUI.Suffix
		}
		if prefix == "" {
			prefix = "⚠ "
		}
	case NotificationError:
		style = ia.errorStyle
		if ia.config != nil {
			prefix = ia.config.General.LogErrorUI.Prefix
			suffix = ia.config.General.LogErrorUI.Suffix
		}
		if prefix == "" {
			prefix = "✗ "
		}
	}

	// Format message with prefix and suffix
	message := prefix + notification.Message + suffix

	// Truncate if needed
	if len(message) > ia.width {
		message = message[:ia.width-3] + "..."
	}

	return style.Render(message)
}