package style

import "github.com/gookit/color"

type Sprinter interface {
	Sprint(a ...interface{}) string
	Sprintf(format string, a ...interface{}) string
}

type TextStyle struct {
	fg         *Color
	bg         *Color
	decoration Decoration

	style Sprinter
}

func New() TextStyle {
	s := TextStyle{}
	s.style = s.deriveStyle()

	return s
}

func (b TextStyle) Sprint(a ...interface{}) string {
	return b.style.Sprint(a...)
}

func (b TextStyle) Sprintf(format string, a ...interface{}) string {
	return b.style.Sprintf(format, a...)
}

func (b TextStyle) SetBold() TextStyle {
	b.decoration.SetBold()
	b.style = b.deriveStyle()

	return b
}

func (b TextStyle) SetUnderline() TextStyle {
	b.decoration.SetUnderline()
	b.style = b.deriveStyle()

	return b
}

func (b TextStyle) SetReverse() TextStyle {
	b.decoration.SetReverse()
	b.style = b.deriveStyle()

	return b
}

func (b TextStyle) SetBg(color Color) TextStyle {
	b.bg = &color
	b.style = b.deriveStyle()

	return b
}

func (b TextStyle) SetFg(color Color) TextStyle {
	b.fg = &color
	b.style = b.deriveStyle()

	return b
}

func (b TextStyle) deriveStyle() Sprinter {
	if b.fg == nil && b.bg == nil {
		return color.Style(b.decoration.ToOpts())
	}

	isRgb := (b.fg != nil && b.fg.IsRGB()) || (b.bg != nil && b.bg.IsRGB())
	if isRgb {
		return b.deriveRGBStyle()
	}

	return b.deriveBasicStyle()
}

func (b TextStyle) deriveBasicStyle() color.Style {
	style := make([]color.Color, 0, 5)

	if b.fg != nil {
		style = append(style, *b.fg.basic)
	}

	if b.bg != nil {
		style = append(style, *b.bg.basic)
	}

	style = append(style, b.decoration.ToOpts()...)

	return color.Style(style)
}

func (b TextStyle) deriveRGBStyle() *color.RGBStyle {
	style := &color.RGBStyle{}

	if b.fg != nil {
		style.SetFg(*b.fg.ToRGB(false).rgb)
	}

	if b.bg != nil {
		// We need to convert the bg firstly to a foreground color,
		// For more info see
		style.SetBg(*b.bg.ToRGB(true).rgb)
	}

	style.SetOpts(b.decoration.ToOpts())

	return style
}
