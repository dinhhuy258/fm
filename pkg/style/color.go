package style

import "github.com/gookit/color"

type Color struct {
	basic *color.Basic
}

func NewBasicColor(cl color.Basic) Color {
	c := Color{}
	c.basic = &cl

	return c
}
