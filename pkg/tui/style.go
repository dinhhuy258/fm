package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/gookit/color"

	"github.com/dinhhuy258/fm/pkg/config"
)

const (
	// Hex color validation constants
	shortHexLength = 4 // #RGB
	longHexLength  = 7 // #RRGGBB
	hexPrefix      = "#"
)

// ColorEntry represents a color with both foreground and background variants
type ColorEntry struct {
	Foreground string
	Background string
}

// colorMap maps color names to their hex values
var colorMap = map[string]ColorEntry{
	"default": {color.FgWhite.RGB().Hex(), color.BgBlack.RGB().Hex()},
	"black":   {color.FgBlack.RGB().Hex(), color.BgBlack.RGB().Hex()},
	"red":     {color.FgRed.RGB().Hex(), color.BgRed.RGB().Hex()},
	"green":   {color.FgGreen.RGB().Hex(), color.BgGreen.RGB().Hex()},
	"yellow":  {color.FgYellow.RGB().Hex(), color.BgYellow.RGB().Hex()},
	"blue":    {color.FgBlue.RGB().Hex(), color.BgBlue.RGB().Hex()},
	"magenta": {color.FgMagenta.RGB().Hex(), color.BgMagenta.RGB().Hex()},
	"cyan":    {color.FgCyan.RGB().Hex(), color.BgCyan.RGB().Hex()},
	"white":   {color.FgWhite.RGB().Hex(), color.BgWhite.RGB().Hex()},
}

// isValidHexValue validates if a string is a valid hex color value.
func isValidHexValue(hex string) bool {
	if hex == "" {
		return false
	}

	if len(hex) != shortHexLength && len(hex) != longHexLength {
		return false
	}

	if !strings.HasPrefix(hex, hexPrefix) {
		return false
	}

	// Validate hex characters
	for _, char := range hex[1:] {
		if !isHexChar(char) {
			return false
		}
	}

	return true
}

// isHexChar checks if a character is a valid hexadecimal digit
func isHexChar(char rune) bool {
	return (char >= '0' && char <= '9') ||
		(char >= 'A' && char <= 'F') ||
		(char >= 'a' && char <= 'f')
}

// parseColor converts a color string to a valid hex color.
func parseColor(colorStr string, isBg bool) string {
	if colorStr == "" {
		return getDefaultColor(isBg)
	}

	// Check if it's a valid hex color
	if isValidHexValue(colorStr) {
		return colorStr
	}

	// Check if it's a named color
	if colorEntry, exists := colorMap[strings.ToLower(colorStr)]; exists {
		if isBg {
			return colorEntry.Background
		}

		return colorEntry.Foreground
	}

	// Fall back to default color
	return getDefaultColor(isBg)
}

// getDefaultColor returns the default foreground or background color
func getDefaultColor(isBg bool) string {
	if isBg {
		return colorMap["default"].Background
	}

	return colorMap["default"].Foreground
}

// fromStyleConfig converts a config.StyleConfig to a lipgloss.Style.
func fromStyleConfig(styleConfig *config.StyleConfig) lipgloss.Style {
	style := lipgloss.NewStyle()

	if styleConfig == nil {
		return style
	}

	// Parse and store colors before applying decorations
	fgColor := parseColor(styleConfig.Fg, false)
	bgColor := parseColor(styleConfig.Bg, true)

	// Set foreground and background colors
	style = style.Foreground(lipgloss.Color(fgColor))
	style = style.Background(lipgloss.Color(bgColor))

	// Apply decorations
	for _, decoration := range styleConfig.Decorations {
		switch strings.ToLower(decoration) {
		case "bold":
			style = style.Bold(true)
		case "italic":
			style = style.Italic(true)
		case "underline", "underscore":
			style = style.Underline(true)
		case "reverse":
			// For reverse, swap the already-parsed foreground and background colors
			style = style.Foreground(lipgloss.Color(bgColor))
			style = style.Background(lipgloss.Color(fgColor))
		}
	}

	return style
}
