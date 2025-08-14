package tui

import (
	"fmt"

	"github.com/dinhhuy258/fm/pkg/config"
)

// ModeManager handles mode definitions, switching, and state management
type ModeManager struct {
	currentMode  string
	customModes  map[string]*config.ModeConfig
	builtinModes map[string]*config.ModeConfig
	modeHistory  []string
}

// NewModeManager creates a new mode manager from the config
func NewModeManager() *ModeManager {
	cfg := config.AppConfig
	mm := &ModeManager{
		currentMode:  "default",
		customModes:  cfg.Modes.Customs,
		builtinModes: cfg.Modes.Builtins,
		modeHistory:  []string{"default"},
	}

	return mm
}

// GetCurrentMode returns the name of the current mode
func (mm *ModeManager) GetCurrentMode() string {
	return mm.currentMode
}

// SwitchToMode switches to the specified mode
func (mm *ModeManager) SwitchToMode(modeName string) error {
	if !mm.modeExists(modeName) {
		return fmt.Errorf("mode '%s' does not exist", modeName)
	}

	// Update mode history
	mm.modeHistory = append(mm.modeHistory, mm.currentMode)
	mm.currentMode = modeName

	return nil
}

// GetPreviousMode returns the previous mode from history
func (mm *ModeManager) GetPreviousMode() string {
	if len(mm.modeHistory) > 1 {
		return mm.modeHistory[len(mm.modeHistory)-2]
	}

	return "default"
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

// GetModeConfig returns the mode configuration for the specified mode
func (mm *ModeManager) GetModeConfig(modeName string) (*config.ModeConfig, error) {
	// Check custom modes first
	if modeConfig, exists := mm.customModes[modeName]; exists {
		return modeConfig, nil
	}

	// Check builtin modes
	if modeConfig, exists := mm.builtinModes[modeName]; exists {
		return modeConfig, nil
	}

	return nil, fmt.Errorf("mode '%s' not found", modeName)
}

// GetAvailableModes returns a list of all available mode names
func (mm *ModeManager) GetAvailableModes() []string {
	var modes []string

	// Add custom modes
	for modeName := range mm.customModes {
		modes = append(modes, modeName)
	}

	// Add builtin modes
	for modeName := range mm.builtinModes {
		modes = append(modes, modeName)
	}

	return modes
}