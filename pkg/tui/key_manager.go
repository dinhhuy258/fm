package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/dinhhuy258/fm/pkg/config"
)

// KeyManager generates and manages keybindings from mode configurations
type KeyManager struct {
	modeManager *ModeManager
}

// NewKeyManager creates a new key manager
func NewKeyManager(modeManager *ModeManager) *KeyManager {
	return &KeyManager{
		modeManager: modeManager,
	}
}

// ResolveKeyAction resolves a tea.KeyMsg to its corresponding action configuration
func (km *KeyManager) ResolveKeyAction(msg tea.KeyMsg) (string, *config.ActionConfig, bool) {
	// Convert the tea.KeyMsg to a key string
	keyStr := km.keyMsgToString(msg)

	if action := km.getActionForKey(keyStr); action != nil {
		return keyStr, action, true
	}

	return "", nil, false
}

// keyMsgToString converts a tea.KeyMsg to a string representation
func (km *KeyManager) keyMsgToString(msg tea.KeyMsg) string {
	// Handle special cases first
	switch msg.Type {
	case tea.KeyCtrlC:
		return "ctrl+c"
	case tea.KeyEsc:
		return "esc"
	case tea.KeyEnter:
		return "enter"
	case tea.KeySpace:
		return "space"
	case tea.KeyTab:
		return "tab"
	case tea.KeyBackspace:
		return "backspace"
	case tea.KeyDelete:
		return "delete"
	case tea.KeyUp:
		return "up"
	case tea.KeyDown:
		return "down"
	case tea.KeyLeft:
		return "left"
	case tea.KeyRight:
		return "right"
	case tea.KeyHome:
		return "home"
	case tea.KeyEnd:
		return "end"
	case tea.KeyPgUp:
		return "pgup"
	case tea.KeyPgDown:
		return "pgdown"
	}

	// Handle special control keys that are already handled above
	// Most control sequences are handled by specific tea.Key* constants

	// Handle regular character keys
	if msg.Type == tea.KeyRunes && len(msg.Runes) > 0 {
		return string(msg.Runes)
	}

	// Fallback to string representation
	return msg.String()
}

// getActionForKey gets the action config for a key in the current mode
func (km *KeyManager) getActionForKey(keyStr string) *config.ActionConfig {
	currentMode := km.modeManager.GetCurrentMode()
	modeConfig, err := km.modeManager.GetModeConfig(currentMode)
	if err != nil {
		return nil
	}

	// Check on_keys first
	if modeConfig.KeyBindings.OnKeys != nil {
		if action, exists := modeConfig.KeyBindings.OnKeys[keyStr]; exists {
			return action
		}
	}

	// Return default action if no specific key binding found
	if modeConfig.KeyBindings.Default != nil {
		return modeConfig.KeyBindings.Default
	}

	return nil
}