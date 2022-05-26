package style

import (
	"github.com/dinhhuy258/gocui"
	"github.com/gookit/color"
)

var stringToColorMap = map[string]color.Color{
	"black":   color.Black,
	"red":     color.Red,
	"green":   color.Green,
	"yellow":  color.Yellow,
	"blue":    color.Blue,
	"magenta": color.Magenta,
	"cyan":    color.Cyan,
	"white":   color.White,
}

var stringToGoCuiColorMap = map[string]gocui.Attribute{
	"black":   gocui.ColorBlack,
	"red":     gocui.ColorRed,
	"green":   gocui.ColorGreen,
	"yellow":  gocui.ColorYellow,
	"blue":    gocui.ColorBlue,
	"magenta": gocui.ColorMagenta,
	"cyan":    gocui.ColorCyan,
	"white":   gocui.ColorWhite,
}

func StringToGoCuiColor(colorStr string) gocui.Attribute {
	c, hasColor := stringToGoCuiColorMap[colorStr]
	if !hasColor {
		c = gocui.ColorWhite
	}

	return c
}

func StringToColor(colorStr string) color.Color {
	c, hasColor := stringToColorMap[colorStr]
	if !hasColor {
		c = color.White
	}

	return c
}

type Color struct {
	basic *color.Basic
}

func NewBasicColor(cl color.Basic) Color {
	c := Color{}
	c.basic = &cl

	return c
}
