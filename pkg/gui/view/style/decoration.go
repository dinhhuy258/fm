package style

import "github.com/gookit/color"

type Decoration struct {
	bold      bool
	underline bool
	reverse   bool
	italic    bool
}

func (d *Decoration) SetBold() {
	d.bold = true
}

func (d *Decoration) SetUnderline() {
	d.underline = true
}

func (d *Decoration) SetReverse() {
	d.reverse = true
}

func (d *Decoration) SetItalic() {
	d.italic = true
}

func (d Decoration) ToOpts() color.Opts {
	opts := make([]color.Color, 0, 3)

	if d.bold {
		opts = append(opts, color.OpBold)
	}

	if d.underline {
		opts = append(opts, color.OpUnderscore)
	}

	if d.reverse {
		opts = append(opts, color.OpReverse)
	}

	if d.italic {
		opts = append(opts, color.OpItalic)
	}

	return opts
}
