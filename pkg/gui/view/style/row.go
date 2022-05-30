package style

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

var ErrInvalidRowData = errors.New("invalid row data")

type Row struct {
	width int
	cells []*cell
}

func (r *Row) SetWidth(width int) {
	r.width = width
}

func (r *Row) AddCell(percentage int, leftAlign bool) {
	r.cells = append(r.cells, &cell{
		percentage: percentage,
		leftAlign:  leftAlign,
	})
}

func (r *Row) Sprint(cells []string) (string, error) {
	if len(cells) != len(r.cells) {
		return "", ErrInvalidRowData
	}

	t := ""

	for i, v := range cells {
		c := r.cells[i]
		w := int(float32(c.percentage) / 100.0 * float32(r.width))
		t += c.sprint(v, w)
	}

	return t, nil
}

type cell struct {
	percentage int
	leftAlign  bool
}

func (c *cell) sprint(t string, w int) string {
	if utf8.RuneCountInString(t) > w {
		t = t[:w-1]
	}

	if c.leftAlign {
		t = paddingRight(t, w, " ")
	} else {
		t = fmt.Sprintf("%"+fmt.Sprintf("%v", w)+"s", t)
	}

	return t
}

// paddingRight right-pads the string with pad up to len runes
func paddingRight(str string, length int, pad string) string {
	return str + strings.Repeat(pad, length-utf8.RuneCountInString(str))
}
