package style

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

var ErrInvalidRowData = errors.New("invalid row data")

type CellValueComponent struct {
	Value string
	Style *TextStyle
}

type CellValue interface{} // FIXME: find out how to get `string| []CellValueComponent`

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

func (r *Row) Sprint(cellVals []CellValue) (string, error) {
	if len(cellVals) != len(r.cells) {
		return "", ErrInvalidRowData
	}

	t := ""

	for i, v := range cellVals {
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

func styleString(val string, style *TextStyle) string {
	if style != nil {
		return style.Sprint(val)
	}

	return val
}

func (c *cell) sprintCellComponents(cellComponents []CellValueComponent, w int) string {
	originalLine := ""
	for _, cellComponent := range cellComponents {
		originalLine += cellComponent.Value
	}

	if utf8.RuneCountInString(originalLine) > w {
		originalLine = originalLine[:w-1]
	}

	originalLineSize := len(originalLine)
	spaceCount := 0

	if c.leftAlign {
		originalLine = paddingRight(originalLine, w, " ")
		spaceCount = len(originalLine) - originalLineSize
	} else {
		originalLine = fmt.Sprintf("%"+fmt.Sprintf("%v", w)+"s", originalLine)
		spaceCount = len(originalLine) - originalLineSize
	}

	originalLineSize = len(originalLine)

	line := ""
	lineSize := 0

	for _, cellComponent := range cellComponents {
		cellVal := cellComponent.Value
		lineSize += len(cellVal)

		if lineSize <= originalLineSize {
			line += styleString(cellVal, cellComponent.Style)
		} else {
			// lineSize > originalLineSize
			discardSize := lineSize - originalLineSize
			cellVal = cellVal[:len(cellVal)-discardSize]
			line += styleString(cellVal, cellComponent.Style)

			break
		}
	}

	if c.leftAlign {
		line += strings.Repeat(" ", spaceCount)
	} else {
		line = strings.Repeat(" ", spaceCount) + line
	}

	return line
}

func (c *cell) sprintString(val string, w int) string {
	line := val

	if utf8.RuneCountInString(line) > w {
		line = line[:w-1]
	}

	if c.leftAlign {
		line = paddingRight(line, w, " ")
	} else {
		line = fmt.Sprintf("%"+fmt.Sprintf("%v", w)+"s", line)
	}

	return line
}

func (c *cell) sprint(cv CellValue, w int) string {
	if w <= 0 {
		return ""
	}

	switch cellValue := cv.(type) {
	case string:
		return c.sprintString(cellValue, w)
	case []CellValueComponent:
		return c.sprintCellComponents(cellValue, w)
	}

	return ""
}

// paddingRight right-pads the string with pad up to len runes
func paddingRight(str string, length int, pad string) string {
	return str + strings.Repeat(pad, length-utf8.RuneCountInString(str))
}
