package tui

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/dinhhuy258/fm/pkg/config"
)

// DynamicKeyMap generates and manages keybindings from mode configurations
type DynamicKeyMap struct {
	modeManager *ModeManager
	keyCache    map[string]map[string]key.Binding // cache keybindings per mode
}

// NewDynamicKeyMap creates a new dynamic keymap
func NewDynamicKeyMap(modeManager *ModeManager) *DynamicKeyMap {
	return &DynamicKeyMap{
		modeManager: modeManager,
		keyCache:    make(map[string]map[string]key.Binding),
	}
}

// GetKeyBindings returns all key bindings for the current mode
func (dk *DynamicKeyMap) GetKeyBindings() map[string]key.Binding {
	currentMode := dk.modeManager.GetCurrentMode()

	// Return cached bindings if available
	if bindings, exists := dk.keyCache[currentMode]; exists {
		return bindings
	}

	// Generate and cache new bindings
	bindings := dk.generateKeyBindings(currentMode)
	dk.keyCache[currentMode] = bindings

	return bindings
}

// MatchesKey checks if a tea.KeyMsg matches any configured key for the current mode
func (dk *DynamicKeyMap) MatchesKey(msg tea.KeyMsg) (string, *config.ActionConfig, bool) {
	// Convert the tea.KeyMsg to a key string
	keyStr := dk.keyMsgToString(msg)

	// Try to get action from mode config
	if action, err := dk.modeManager.GetKeyBinding(keyStr); err == nil {
		return keyStr, action, true
	}

	// Check for alternative key representations
	altKeys := dk.getAlternativeKeys(keyStr)
	for _, altKey := range altKeys {
		if action, err := dk.modeManager.GetKeyBinding(altKey); err == nil {
			return altKey, action, true
		}
	}

	return "", nil, false
}

// generateKeyBindings generates key.Binding objects for a specific mode
func (dk *DynamicKeyMap) generateKeyBindings(modeName string) map[string]key.Binding {
	bindings := make(map[string]key.Binding)

	modeConfig, err := dk.modeManager.GetModeConfig(modeName)
	if err != nil {
		// Return empty bindings for modes without config (like default)
		return bindings
	}

	// Generate bindings for each key in on_keys
	if modeConfig.KeyBindings.OnKeys != nil {
		for keyStr, action := range modeConfig.KeyBindings.OnKeys {
			binding := dk.createKeyBinding(keyStr, action)
			bindings[keyStr] = binding
		}
	}

	return bindings
}

// createKeyBinding creates a key.Binding from a key string and action config
func (dk *DynamicKeyMap) createKeyBinding(keyStr string, action *config.ActionConfig) key.Binding {
	keys := dk.parseKeyString(keyStr)
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
func (dk *DynamicKeyMap) parseKeyString(keyStr string) []string {
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
func (dk *DynamicKeyMap) keyMsgToString(msg tea.KeyMsg) string {
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
func (dk *DynamicKeyMap) getAlternativeKeys(keyStr string) []string {
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

// ClearCache clears the keymap cache (useful when modes change)
func (dk *DynamicKeyMap) ClearCache() {
	dk.keyCache = make(map[string]map[string]key.Binding)
}

// HasBinding checks if a key string has a binding in the current mode
func (dk *DynamicKeyMap) HasBinding(keyStr string) bool {
	return dk.modeManager.HasKeyBinding(keyStr)
}

// GetBinding returns the key.Binding for a specific key string in current mode
func (dk *DynamicKeyMap) GetBinding(keyStr string) (key.Binding, bool) {
	bindings := dk.GetKeyBindings()
	if binding, exists := bindings[keyStr]; exists {
		return binding, true
	}

	return key.Binding{}, false
}

// GetAllBindingsForHelp returns all bindings formatted for help display
func (dk *DynamicKeyMap) GetAllBindingsForHelp() []key.Binding {
	bindings := dk.GetKeyBindings()
	var result []key.Binding

	for _, binding := range bindings {
		result = append(result, binding)
	}

	return result
}
