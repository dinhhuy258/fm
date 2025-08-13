package tui

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/dinhhuy258/fm/pkg/config"
)

// KeyManager generates and manages keybindings from mode configurations
type KeyManager struct {
	modeManager *ModeManager
	keyBindings map[string]map[string]key.Binding
}

// NewKeyManager creates a new key manager and preloads all keybindings
func NewKeyManager(modeManager *ModeManager) *KeyManager {
	km := &KeyManager{
		modeManager: modeManager,
		keyBindings: make(map[string]map[string]key.Binding),
	}

	// Preload keybindings for all available modes
	km.loadAllKeyBindings()

	return km
}

// MatchesKey checks if a tea.KeyMsg matches any configured key for the current mode
func (km *KeyManager) MatchesKey(msg tea.KeyMsg) (string, *config.ActionConfig, bool) {
	// Convert the tea.KeyMsg to a key string
	keyStr := km.keyMsgToString(msg)

	// Try to get action from current mode's key bindings
	if action := km.getActionForKey(keyStr); action != nil {
		return keyStr, action, true
	}

	// Check for alternative key representations
	altKeys := km.getAlternativeKeys(keyStr)
	for _, altKey := range altKeys {
		if action := km.getActionForKey(altKey); action != nil {
			return altKey, action, true
		}
	}

	return "", nil, false
}

// generateKeyBindings generates key.Binding objects for a specific mode
func (km *KeyManager) generateKeyBindings(modeName string) map[string]key.Binding {
	bindings := make(map[string]key.Binding)

	modeConfig, err := km.modeManager.GetModeConfig(modeName)
	if err != nil {
		// Return empty bindings for modes without config (like default)
		return bindings
	}

	// Generate bindings for each key in on_keys
	if modeConfig.KeyBindings.OnKeys != nil {
		for keyStr, action := range modeConfig.KeyBindings.OnKeys {
			binding := km.createKeyBinding(keyStr, action)
			bindings[keyStr] = binding
		}
	}

	return bindings
}

// createKeyBinding creates a key.Binding from a key string and action config
func (km *KeyManager) createKeyBinding(keyStr string, action *config.ActionConfig) key.Binding {
	keys := km.parseKeyString(keyStr)
	help := action.Help
	if help == "" {
		help = keyStr
	}

	return key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(keyStr, help),
	)
}

// parseKeyString converts a config key string to bubbletea key strings
func (km *KeyManager) parseKeyString(keyStr string) []string {
	// Handle special key mappings
	keyMappings := map[string]string{
		"esc":       "esc",
		"enter":     "enter",
		"space":     " ",
		"tab":       "tab",
		"backspace": "backspace",
		"delete":    "delete",
		"up":        "up",
		"down":      "down",
		"left":      "left",
		"right":     "right",
		"home":      "home",
		"end":       "end",
		"pgup":      "pgup",
		"pgdown":    "pgdown",
	}

	// Check if it's a mapped special key
	if mapped, exists := keyMappings[strings.ToLower(keyStr)]; exists {
		return []string{mapped}
	}

	// Handle control sequences like "ctrl+c"
	if strings.Contains(keyStr, "ctrl+") {
		return []string{keyStr}
	}

	// Handle alt sequences like "alt+h"
	if strings.Contains(keyStr, "alt+") {
		return []string{keyStr}
	}

	// Handle shift sequences like "shift+tab"
	if strings.Contains(keyStr, "shift+") {
		return []string{keyStr}
	}

	// Single character keys
	return []string{keyStr}
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

// getAlternativeKeys returns alternative representations of a key
func (km *KeyManager) getAlternativeKeys(keyStr string) []string {
	alternatives := []string{}

	// Map alternative representations
	altMappings := map[string][]string{
		" ":      {"space"},
		"space":  {" "},
		"enter":  {"return"},
		"return": {"enter"},
		"esc":    {"escape"},
		"escape": {"esc"},
		"ctrl+c": {"q"}, // Common quit alternative
	}

	if alts, exists := altMappings[keyStr]; exists {
		alternatives = append(alternatives, alts...)
	}

	// Handle case variations for control sequences
	if strings.Contains(keyStr, "ctrl+") {
		// Try uppercase version
		re := regexp.MustCompile(`ctrl\+([a-z])`)
		if matches := re.FindStringSubmatch(keyStr); len(matches) > 1 {
			upperKey := "ctrl+" + strings.ToUpper(matches[1])
			alternatives = append(alternatives, upperKey)
		}
	}

	return alternatives
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

// loadAllKeyBindings preloads keybindings for all available modes
func (km *KeyManager) loadAllKeyBindings() {
	// Get all available modes from the mode manager
	availableModes := km.modeManager.GetAvailableModes()

	// Generate bindings for each mode
	for _, modeName := range availableModes {
		bindings := km.generateKeyBindings(modeName)
		km.keyBindings[modeName] = bindings
	}
}
