package view

import (
	"github.com/dinhhuy258/gocui"
	"github.com/gookit/color"
)

func toGocuiAttribute(c color.Color) gocui.Attribute {
	switch {
	case c == color.Black:
		return gocui.ColorBlack
	case c == color.Red:
		return gocui.ColorRed
	case c == color.Green:
		return gocui.ColorGreen
	case c == color.Yellow:
		return gocui.ColorYellow
	case c == color.Blue:
		return gocui.ColorBlue
	case c == color.Magenta:
		return gocui.ColorMagenta
	case c == color.Cyan:
		return gocui.ColorCyan
	case c == color.White:
		return gocui.ColorWhite
	default:
		return gocui.ColorWhite
	}
}

func toColor(c string) color.Color {
	switch c {
	case "black":
		return color.Black
	case "red":
		return color.Red
	case "green":
		return color.Green
	case "yellow":
		return color.Yellow
	case "blue":
		return color.Blue
	case "magenta":
		return color.Magenta
	case "cyan":
		return color.Cyan
	case "white":
		return color.White
	default:
		return color.White
	}
}
