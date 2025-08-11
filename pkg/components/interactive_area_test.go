package components

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestInteractiveAreaModeSwitch(t *testing.T) {
	ia := NewInteractiveArea()
	ia.SetSize(80, 1)

	// Initial state should be notification mode
	if ia.IsInputMode() {
		t.Error("InteractiveArea should start in notification mode")
	}

	// Show input should switch to input mode
	ia.ShowInput("> ", "test")
	if !ia.IsInputMode() {
		t.Error("ShowInput should switch to input mode")
	}

	// Hide input should switch back to notification mode
	ia.HideInput()
	if ia.IsInputMode() {
		t.Error("HideInput should switch back to notification mode")
	}
}

func TestInteractiveAreaInputCompletion(t *testing.T) {
	ia := NewInteractiveArea()
	ia.SetSize(80, 1)

	// Switch to input mode
	ia.ShowInput("> ", "test input")

	// Simulate Enter key press
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedIA, cmd := ia.Update(enterMsg)

	// Should be back in notification mode
	if updatedIA.IsInputMode() {
		t.Error("Enter key should switch back to notification mode")
	}

	// Should return InputCompletedMessage
	if cmd == nil {
		t.Error("Enter key should return InputCompletedMessage command")
	} else {
		msg := cmd()
		if completedMsg, ok := msg.(InputCompletedMessage); ok {
			if completedMsg.Value != "test input" {
				t.Errorf("Expected 'test input', got '%s'", completedMsg.Value)
			}
		} else {
			t.Error("Expected InputCompletedMessage")
		}
	}
}

func TestInteractiveAreaInputCancellation(t *testing.T) {
	ia := NewInteractiveArea()
	ia.SetSize(80, 1)

	// Switch to input mode
	ia.ShowInput("> ", "test input")

	// Simulate Escape key press
	escapeMsg := tea.KeyMsg{Type: tea.KeyEsc}
	updatedIA, cmd := ia.Update(escapeMsg)

	// Should be back in notification mode
	if updatedIA.IsInputMode() {
		t.Error("Escape key should switch back to notification mode")
	}

	// Should not return a completion command
	if cmd != nil {
		t.Error("Escape key should not return a completion command")
	}
}

func TestInteractiveAreaNotificationQueuing(t *testing.T) {
	ia := NewInteractiveArea()
	ia.SetSize(80, 1)

	// Switch to input mode
	ia.ShowInput("> ", "test")

	// Try to show a notification while in input mode - should be queued
	cmd := ia.ShowSuccess("This should be queued")
	if cmd != nil {
		t.Error("Notification commands should not be returned while in input mode")
	}

	// Should not have active notification but should have queued notification
	if ia.HasActiveNotification() {
		t.Error("Should not have active notification while in input mode")
	}

	if !ia.HasPendingNotification() {
		t.Error("Should have pending notification while in input mode")
	}

	// Switch back to notification mode - should activate queued notification
	hideCmd := ia.HideInput()

	// Should now have active notification and no queued notification
	if !ia.HasActiveNotification() {
		t.Error("Should have active notification after exiting input mode")
	}

	if ia.HasPendingNotification() {
		t.Error("Should not have pending notification after exiting input mode")
	}

	// Should return auto-clear command for success notification
	if hideCmd == nil {
		t.Error("HideInput should return auto-clear command for success notification")
	}
}

func TestInteractiveAreaMultipleQueuedNotifications(t *testing.T) {
	ia := NewInteractiveArea()
	ia.SetSize(80, 1)

	// Switch to input mode
	ia.ShowInput("> ", "test")

	// Try to show multiple notifications - only the last should be kept
	ia.ShowInfo("First notification")
	ia.ShowWarning("Second notification")
	ia.ShowSuccess("Third notification")

	// Should have pending notification (the last one)
	if !ia.HasPendingNotification() {
		t.Error("Should have pending notification")
	}

	// Exit input mode
	ia.HideInput()

	// Should show the last notification (success)
	if !ia.HasActiveNotification() {
		t.Error("Should have active notification")
	}

	// Verify it's the success message by checking the view contains success icon
	view := ia.View()
	if view == "" {
		t.Error("Should render notification view")
	}
	// Success notifications should have ✓ prefix
	if !strings.Contains(view, "✓") {
		t.Error("Should show success notification (with ✓ prefix)")
	}
}

func TestInteractiveAreaView(t *testing.T) {
	ia := NewInteractiveArea()
	ia.SetSize(80, 1)

	// Test notification view
	ia.ShowInfo("Test notification")
	view := ia.View()
	if view == "" {
		t.Error("Should render notification view")
	}

	// Test input view
	ia.ShowInput("> ", "test input")
	view = ia.View()
	if view == "" {
		t.Error("Should render input view")
	}

	// Test empty view
	ia.HideInput()
	ia.ClearNotification()
	view = ia.View()
	if view != "" {
		t.Error("Should render empty view when no notification or input")
	}
}
