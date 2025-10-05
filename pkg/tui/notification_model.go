package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/dinhhuy258/fm/pkg/config"
)

const autoClearNotificationDuration = 5 * time.Second

// NotificationType represents different types of status notifications
type NotificationType int8

const (
	NotificationSuccess NotificationType = iota
	NotificationInfo
	NotificationWarning
	NotificationError
)

// Notification represents a single status notification with timestamp
type Notification struct {
	Type      NotificationType
	Message   string
	CreatedAt time.Time
}

// NotificationStyles holds cached styles for notifications
type NotificationStyles struct {
	successStyle lipgloss.Style
	infoStyle    lipgloss.Style
	warningStyle lipgloss.Style
	errorStyle   lipgloss.Style
}

// NotificationModel handles notification display and management
type NotificationModel struct {
	width  int
	height int

	activeNotification *Notification
	styles             *NotificationStyles
	isVisible          bool
}

// NewNotificationModel creates a new notification model
func NewNotificationModel() *NotificationModel {
	notificationStyle := &NotificationStyles{
		successStyle: fromStyleConfig(config.AppConfig.General.LogInfoUI.Style),
		infoStyle:    fromStyleConfig(config.AppConfig.General.LogInfoUI.Style),
		warningStyle: fromStyleConfig(config.AppConfig.General.LogWarningUI.Style),
		errorStyle:   fromStyleConfig(config.AppConfig.General.LogErrorUI.Style),
	}

	return &NotificationModel{
		activeNotification: nil,
		isVisible:          true, // Default to visible
		styles:             notificationStyle,
	}
}

// SetSize updates the model dimensions
func (m *NotificationModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// GetSize returns the current dimensions
func (m *NotificationModel) GetSize() (int, int) {
	return m.width, m.height
}

// Show makes the notification visible
func (m *NotificationModel) Show() {
	m.isVisible = true
}

// Hide makes the notification invisible
func (m *NotificationModel) Hide() {
	m.isVisible = false
}

// IsVisible returns whether the notification is currently visible
func (m *NotificationModel) IsVisible() bool {
	return m.isVisible
}

// ShowNotification displays a notification
func (m *NotificationModel) ShowNotification(notificationType NotificationType, message string) tea.Cmd {
	notification := &Notification{
		Type:      notificationType,
		Message:   message,
		CreatedAt: time.Now(),
	}

	m.activeNotification = notification

	// Auto-clear notifications
	return tea.Tick(autoClearNotificationDuration, func(t time.Time) tea.Msg {
		return AutoClearMessage{}
	})
}

// GetActiveNotification returns the current active notification
func (m *NotificationModel) GetActiveNotification() *Notification {
	return m.activeNotification
}

// ClearNotification clears the active notification
func (m *NotificationModel) ClearNotification() {
	m.activeNotification = nil
}

// AutoClearMessage represents a message to auto-clear the notification
type AutoClearMessage struct{}

// View renders the notification view
func (m *NotificationModel) View() string {
	if !m.isVisible || m.activeNotification == nil || m.width <= 0 || m.height <= 0 {
		return ""
	}

	var notificationStyle lipgloss.Style
	var prefix, suffix string
	cfg := config.AppConfig.General

	// Get style and prefix/suffix based on notification type
	switch m.activeNotification.Type {
	case NotificationSuccess:
		notificationStyle = m.styles.successStyle
		prefix = cfg.LogInfoUI.Prefix
		suffix = cfg.LogInfoUI.Suffix
	case NotificationInfo:
		notificationStyle = m.styles.infoStyle
		prefix = cfg.LogInfoUI.Prefix
		suffix = cfg.LogInfoUI.Suffix
	case NotificationWarning:
		notificationStyle = m.styles.warningStyle
		prefix = cfg.LogWarningUI.Prefix
		suffix = cfg.LogWarningUI.Suffix
	case NotificationError:
		notificationStyle = m.styles.errorStyle
		prefix = cfg.LogErrorUI.Prefix
		suffix = cfg.LogErrorUI.Suffix
	}

	// Format message with prefix and suffix
	message := prefix + m.activeNotification.Message + suffix

	// Truncate the message if it's wider than the view.
	message = Truncate(message, m.width, "...")

	return notificationStyle.Render(message)
}
