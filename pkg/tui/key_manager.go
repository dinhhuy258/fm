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
func (km *KeyManager) ResolveKeyAction(msg tea.KeyMsg) *config.ActionConfig {
	keyStr := msg.String()

	currentMode := km.modeManager.GetCurrentMode()
	modeConfig := km.modeManager.GetModeConfig(currentMode)
	if modeConfig == nil {
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
