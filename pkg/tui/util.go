package tui

import (
	"strings"

	"github.com/rivo/uniseg"
)

// Truncate shortens a string to a maximum width, adding a suffix if it's too long.
// It respects Unicode grapheme clusters to avoid breaking multi-byte characters.
func Truncate(str string, maxWidth int, suffix string) string {
	strWidth := uniseg.StringWidth(str)
	if strWidth <= maxWidth {
		return str
	}

	suffixWidth := uniseg.StringWidth(suffix)
	if maxWidth <= suffixWidth {
		// Not enough space for the suffix, so we truncate without it.
		// This is an edge case, but we handle it to avoid panics.
		var builder strings.Builder
		truncatedWidth := 0
		gr := uniseg.NewGraphemes(str)
		for gr.Next() {
			runeText := gr.Str()
			runeWidth := uniseg.StringWidth(runeText)
			if truncatedWidth+runeWidth > maxWidth {
				break
			}
			builder.WriteString(runeText)
			truncatedWidth += runeWidth
		}

		return builder.String()
	}

	var builder strings.Builder
	truncatedWidth := 0
	// Use grapheme clusters to avoid breaking multi-byte characters
	gr := uniseg.NewGraphemes(str)
	// Calculate the target width for the main string part
	targetWidth := maxWidth - suffixWidth
	for gr.Next() {
		runeText := gr.Str()
		runeWidth := uniseg.StringWidth(runeText)
		if truncatedWidth+runeWidth > targetWidth {
			break
		}
		builder.WriteString(runeText)
		truncatedWidth += runeWidth
	}

	return builder.String() + suffix
}