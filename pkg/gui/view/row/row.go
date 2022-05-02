package row

import (
	"errors"
	"fmt"
	"github.com/dinhhuy258/fm/pkg/gui/view/style"
)

var ErrInvalidRowData = errors.New("invalid row data")

type Row struct {
	width int
	cells []*rowCell
}

func (r *Row) SetWidth(width int) {
	r.width = width
}

func (r *Row) AddCell(percentage int, leftAlign bool, textStyle *style.TextStyle) {
	r.cells = append(r.cells, &rowCell{
		percentage: percentage,
		leftAlign:  leftAlign,
		textStyle:  textStyle,
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

type rowCell struct {
	percentage int
	leftAlign  bool
	textStyle  *style.TextStyle
}

func (rc *rowCell) sprint(t string, w int) string {
	if rc.leftAlign {
		t = Right(t, w, " ")
	} else {
		t = fmt.Sprintf("%"+fmt.Sprintf("%v", w)+"s", t)
	}

	if rc.textStyle != nil {
		return rc.textStyle.Sprint(t)
	}

	return t
}
