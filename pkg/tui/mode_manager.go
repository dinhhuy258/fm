package tui

import (
	"github.com/dinhhuy258/fm/pkg/config"
)

// ModeManager handles mode definitions, switching, and state management
type ModeManager struct {
	currentMode  string
	customModes  map[string]*config.ModeConfig
	builtinModes map[string]*config.ModeConfig
}

// NewModeManager creates a new mode manager from the config
func NewModeManager() *ModeManager {
	cfg := config.AppConfig
	mm := &ModeManager{
		currentMode:  "default",
		customModes:  cfg.Modes.Customs,
		builtinModes: cfg.Modes.Builtins,
	}

	return mm
}

// GetCurrentMode returns the name of the current mode
func (mm *ModeManager) GetCurrentMode() string {
	return mm.currentMode
}

// SwitchToMode switches to the specified mode
func (mm *ModeManager) SwitchToMode(modeName string) {
	if !mm.modeExists(modeName) {
		return
	}

	mm.currentMode = modeName
}

// GetModeConfig returns the mode configuration for the specified mode
func (mm *ModeManager) GetModeConfig(modeName string) *config.ModeConfig {
	// Check custom modes first
	if modeConfig, exists := mm.customModes[modeName]; exists {
		return modeConfig
	}

	// Check builtin modes
	if modeConfig, exists := mm.builtinModes[modeName]; exists {
		return modeConfig
	}

	return nil
}

// modeExists checks if a mode exists in custom or builtin modes
func (mm *ModeManager) modeExists(modeName string) bool {
	if _, exists := mm.customModes[modeName]; exists {
		return true
	}

	if _, exists := mm.builtinModes[modeName]; exists {
		return true
	}

	return false
}