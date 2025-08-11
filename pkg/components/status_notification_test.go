package components

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestStatusNotificationAutoClear(t *testing.T) {
	notification := NewStatusNotificationView()
	notification.SetSize(80, 1)

	// Test that ShowSuccess returns a command for auto-clearing
	cmd := notification.ShowSuccess("Test success message")
	if cmd == nil {
		t.Error("ShowSuccess should return a tea.Cmd for auto-clear")
	}

	// Verify notification is active
	if !notification.HasActiveNotification() {
		t.Error("Notification should be active after ShowSuccess")
	}

	// Execute the command to get the message
	if cmd != nil {
		// This simulates waiting for the tick
		go func() {
			time.Sleep(100 * time.Millisecond) // Short delay for test
			msg := cmd()
			if _, ok := msg.(AutoClearMessage); !ok {
				t.Error("Command should return AutoClearMessage")
			}
		}()
	}

	// Test manual clearing
	notification.Clear()
	if notification.HasActiveNotification() {
		t.Error("Notification should not be active after Clear()")
	}
}

func TestStatusNotificationTypes(t *testing.T) {
	notification := NewStatusNotificationView()
	notification.SetSize(80, 1)

	tests := []struct {
		name       string
		method     func(string) tea.Cmd
		message    string
		expectsCmd bool
	}{
		{"Success", notification.ShowSuccess, "Success message", true},
		{"Info", notification.ShowInfo, "Info message", false},
		{"Warning", notification.ShowWarning, "Warning message", false},
		{"Error", notification.ShowError, "Error message", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.method(tt.message)
			if tt.expectsCmd && cmd == nil {
				t.Errorf("%s should return a command for auto-clear", tt.name)
			}
			if !tt.expectsCmd && cmd != nil {
				t.Errorf("%s should not return a command", tt.name)
			}

			if !notification.HasActiveNotification() {
				t.Errorf("Notification should be active after %s", tt.name)
			}

			view := notification.View()
			if view == "" {
				t.Errorf("View should not be empty after %s", tt.name)
			}

			notification.Clear()
		})
	}
}

func TestAutoComplexarMessageHandling(t *testing.T) {
	notification := NewStatusNotificationView()
	notification.SetSize(80, 1)

	// Set up a success notification
	notification.ShowSuccess("Test message")

	// Simulate receiving an AutoClearMessage
	updatedNotification, cmd := notification.Update(AutoClearMessage{})
	if cmd != nil {
		t.Error("Update with AutoClearMessage should not return additional commands")
	}

	if updatedNotification.HasActiveNotification() {
		t.Error("Notification should be cleared after AutoClearMessage")
	}
}
