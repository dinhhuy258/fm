package tui

import (
	"fmt"

	"github.com/dinhhuy258/fm/pkg/config"
)

// ModeManager handles mode definitions, switching, and state management
type ModeManager struct {
	currentMode  string
	defaultMode  string
	customModes  map[string]*config.ModeConfig
	builtinModes map[string]*config.ModeConfig
	modeHistory  []string
}

// NewModeManager creates a new mode manager from the config
func NewModeManager(cfg *config.Config) *ModeManager {
	mm := &ModeManager{
		defaultMode:  "default",
		currentMode:  "default",
		customModes:  make(map[string]*config.ModeConfig),
		builtinModes: make(map[string]*config.ModeConfig),
		modeHistory:  []string{"default"},
	}

	// Load modes from config if available
	if cfg != nil && cfg.Modes != nil {
		if cfg.Modes.Customs != nil {
			mm.customModes = cfg.Modes.Customs
		}
		if cfg.Modes.Builtins != nil {
			mm.builtinModes = cfg.Modes.Builtins
		}
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

	// Keep only the last 10 modes in history
	if len(mm.modeHistory) > 10 {
		mm.modeHistory = mm.modeHistory[len(mm.modeHistory)-10:]
	}

	mm.currentMode = modeName

	return nil
}

// GetPreviousMode returns the previous mode from history
func (mm *ModeManager) GetPreviousMode() string {
	if len(mm.modeHistory) > 1 {
		return mm.modeHistory[len(mm.modeHistory)-2]
	}

	return mm.defaultMode
}

// modeExists checks if a mode exists in custom or builtin modes
func (mm *ModeManager) modeExists(modeName string) bool {
	// Check custom modes first
	if _, exists := mm.customModes[modeName]; exists {
		return true
	}

	// Check builtin modes
	if _, exists := mm.builtinModes[modeName]; exists {
		return true
	}

	// Always allow the default mode
	return modeName == mm.defaultMode
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

	// Add default mode
	modes = append(modes, mm.defaultMode)

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
