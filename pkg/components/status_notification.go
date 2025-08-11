package components

import (
	"strings"
	"time"

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

// StatusNotification represents a single status notification with auto-clear functionality
type StatusNotification struct {
	Type      NotificationType
	Message   string
	CreatedAt time.Time
}

// StatusNotificationView provides a status notification display with auto-clear
type StatusNotificationView struct {
	width  int
	height int

	// Current active notification
	activeNotification *StatusNotification
	clearTimer         *time.Timer

	// Styling based on config
	successStyle lipgloss.Style
	infoStyle    lipgloss.Style
	warningStyle lipgloss.Style
	errorStyle   lipgloss.Style

	// Configuration
	config *config.Config
}

// AutoClearMessage represents a message to auto-clear the notification
type AutoClearMessage struct{}

// NewStatusNotificationView creates a new status notification view
func NewStatusNotificationView() *StatusNotificationView {
	cfg := config.AppConfig

	notification := &StatusNotificationView{
		activeNotification: nil,
		clearTimer:         nil,
		config:             cfg,
	}

	notification.initStyles()

	return notification
}

// initStyles initializes lipgloss styles from config
func (snv *StatusNotificationView) initStyles() {
	// Default styles with good visibility
	snv.successStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")).
		Bold(true)

	snv.infoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(false)

	snv.warningStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFF00")).
		Bold(true)

	snv.errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000")).
		Bold(true)

	// Override with config if available (reusing existing log UI configs)
	if snv.config != nil && snv.config.General.LogInfoUI.Style != nil {
		snv.infoStyle = convertConfigStyleToLipgloss(snv.config.General.LogInfoUI.Style)
		snv.successStyle = convertConfigStyleToLipgloss(
			snv.config.General.LogInfoUI.Style,
		).Foreground(lipgloss.Color("#00FF00"))
	}
	if snv.config != nil && snv.config.General.LogWarningUI.Style != nil {
		snv.warningStyle = convertConfigStyleToLipgloss(snv.config.General.LogWarningUI.Style)
	}
	if snv.config != nil && snv.config.General.LogErrorUI.Style != nil {
		snv.errorStyle = convertConfigStyleToLipgloss(snv.config.General.LogErrorUI.Style)
	}
}

// SetSize updates the notification view size
func (snv *StatusNotificationView) SetSize(width, height int) {
	snv.width = width
	snv.height = height
}

// ShowNotification displays a notification with optional auto-clear
func (snv *StatusNotificationView) ShowNotification(notificationType NotificationType, message string) tea.Cmd {
	// Set the new notification
	snv.activeNotification = &StatusNotification{
		Type:      notificationType,
		Message:   message,
		CreatedAt: time.Now(),
	}

	// Auto-clear success notifications after 5 seconds
	if notificationType == NotificationSuccess {
		return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
			return AutoClearMessage{}
		})
	}

	return nil
}

// ShowSuccess displays a success notification (auto-clears in 5 seconds)
func (snv *StatusNotificationView) ShowSuccess(message string) tea.Cmd {
	return snv.ShowNotification(NotificationSuccess, message)
}

// ShowInfo displays an info notification (persists until manually cleared)
func (snv *StatusNotificationView) ShowInfo(message string) tea.Cmd {
	return snv.ShowNotification(NotificationInfo, message)
}

// ShowWarning displays a warning notification (persists until manually cleared)
func (snv *StatusNotificationView) ShowWarning(message string) tea.Cmd {
	return snv.ShowNotification(NotificationWarning, message)
}

// ShowError displays an error notification (persists until manually cleared)
func (snv *StatusNotificationView) ShowError(message string) tea.Cmd {
	return snv.ShowNotification(NotificationError, message)
}

// Clear clears the current notification
func (snv *StatusNotificationView) Clear() {
	if snv.clearTimer != nil {
		snv.clearTimer.Stop()
		snv.clearTimer = nil
	}
	snv.activeNotification = nil
}

// HasActiveNotification returns whether there's an active notification
func (snv *StatusNotificationView) HasActiveNotification() bool {
	return snv.activeNotification != nil
}

// Update handles input messages
func (snv *StatusNotificationView) Update(msg tea.Msg) (*StatusNotificationView, tea.Cmd) {
	switch msg := msg.(type) {
	case AutoClearMessage:
		// Auto-clear the notification
		snv.Clear()

		return snv, nil
	case tea.KeyMsg:
		// Allow users to manually clear notifications with Escape
		if msg.String() == "esc" && snv.HasActiveNotification() {
			snv.Clear()

			return snv, nil
		}
	}

	return snv, nil
}

// View renders the status notification view
func (snv *StatusNotificationView) View() string {
	if snv.width <= 0 || snv.height <= 0 || snv.activeNotification == nil {
		return ""
	}

	notification := snv.activeNotification
	var style lipgloss.Style
	var prefix, suffix string

	// Get style and prefix/suffix based on notification type
	switch notification.Type {
	case NotificationSuccess:
		style = snv.successStyle
		if snv.config != nil {
			prefix = snv.config.General.LogInfoUI.Prefix
			suffix = snv.config.General.LogInfoUI.Suffix
		}
		if prefix == "" {
			prefix = "✓ "
		}
	case NotificationInfo:
		style = snv.infoStyle
		if snv.config != nil {
			prefix = snv.config.General.LogInfoUI.Prefix
			suffix = snv.config.General.LogInfoUI.Suffix
		}
		if prefix == "" {
			prefix = "ℹ "
		}
	case NotificationWarning:
		style = snv.warningStyle
		if snv.config != nil {
			prefix = snv.config.General.LogWarningUI.Prefix
			suffix = snv.config.General.LogWarningUI.Suffix
		}
		if prefix == "" {
			prefix = "⚠ "
		}
	case NotificationError:
		style = snv.errorStyle
		if snv.config != nil {
			prefix = snv.config.General.LogErrorUI.Prefix
			suffix = snv.config.General.LogErrorUI.Suffix
		}
		if prefix == "" {
			prefix = "✗ "
		}
	}

	// Format message with prefix and suffix
	message := prefix + notification.Message + suffix

	// Truncate if needed
	if len(message) > snv.width {
		message = message[:snv.width-3] + "..."
	}

	// Apply styling and center the message
	styledMessage := style.Render(message)

	// Create lines to fill the height
	lines := make([]string, snv.height)

	// Place the notification message in the first line
	lines[0] = styledMessage

	// Fill remaining lines with empty strings
	for i := 1; i < snv.height; i++ {
		lines[i] = ""
	}

	return strings.Join(lines, "\n")
}
