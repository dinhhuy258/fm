package style

import "github.com/gookit/color"

type Decoration struct {
	bold      bool
	underscore bool
	reverse   bool
	italic    bool
}

func (d *Decoration) SetBold() {
	d.bold = true
}

func (d *Decoration) SetUnderscore() {
	d.underscore = true
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

	if d.underscore {
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
