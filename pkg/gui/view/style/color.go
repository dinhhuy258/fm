package style

import (
	"github.com/gookit/color"
)

var ColorMap = map[string]struct {
	Foreground color.Color
	Background color.Color
}{
	"default": {color.FgWhite, color.BgBlack},
	"black":   {color.FgBlack, color.BgBlack},
	"red":     {color.FgRed, color.BgRed},
	"green":   {color.FgGreen, color.BgGreen},
	"yellow":  {color.FgYellow, color.BgYellow},
	"blue":    {color.FgBlue, color.BgBlue},
	"magenta": {color.FgMagenta, color.BgMagenta},
	"cyan":    {color.FgCyan, color.BgCyan},
	"white":   {color.FgWhite, color.BgWhite},
}

type Color struct {
	rgb   *color.RGBColor
	basic *color.Color
}

func NewRGBColor(cl color.RGBColor) Color {
	c := Color{}
	c.rgb = &cl

	return c
}

func NewBasicColor(cl color.Color) Color {
	c := Color{}
	c.basic = &cl
	return c
}

func (c Color) IsRGB() bool {
	return c.rgb != nil
}

func (c Color) ToRGB(isBg bool) Color {
	if c.IsRGB() {
		return c
	}

	if isBg {
		// We need to convert bg color to fg color
		// This is a gookit/color bug,
		// https://github.com/gookit/color/issues/39
		return NewRGBColor((*c.basic - 10).RGB())
	}

	return NewRGBColor(c.basic.RGB())
}

func getColor(val string, isBg bool) Color {
	if isValidHexValue(val) {
		return NewRGBColor(color.HEX(val, isBg))
	}

	color, hasValue := ColorMap[val]
	if !hasValue {
		color, _ = ColorMap["default"]
	}

	if isBg {
		return NewBasicColor(color.Background)
	}

	return NewBasicColor(color.Foreground)
}

func isValidHexValue(v string) bool {
	if len(v) != 4 && len(v) != 7 {
		return false
	}

	if v[0] != '#' {
		return false
	}

	for _, char := range v[1:] {
		switch char {
		case '0',
			'1',
			'2',
			'3',
			'4',
			'5',
			'6',
			'7',
			'8',
			'9',
			'a',
			'b',
			'c',
			'd',
			'e',
			'f',
			'A',
			'B',
			'C',
			'D',
			'E',
			'F':
			continue
		default:
			return false
		}
	}

	return true
}
