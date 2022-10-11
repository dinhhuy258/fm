package view

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dinhhuy258/fm/pkg/gui/view/style"
	"github.com/rivo/uniseg"
)

var ErrInvalidRowData = errors.New("invalid row data")

type ColumnValueComponent struct {
	Value string
	Style *style.TextStyle
}

type ColumnValue interface{} // FIXME: find out how to get `string| []ColumnValueComponent`

type Row struct {
	width   int
	columns []*column
}

func (r *Row) SetWidth(width int) {
	r.width = width
}

func (r *Row) AddColumn(percentage int, leftAlign bool) {
	r.columns = append(r.columns, &column{
		percentage: percentage,
		leftAlign:  leftAlign,
	})
}

func (r *Row) Sprint(colVals []ColumnValue) (string, error) {
	if len(colVals) != len(r.columns) {
		return "", ErrInvalidRowData
	}

	t := ""

	for i, v := range colVals {
		c := r.columns[i]
		w := int(float32(c.percentage) / 100.0 * float32(r.width))
		t += c.sprint(v, w)
	}

	return t, nil
}

type column struct {
	percentage int
	leftAlign  bool
}

func styleString(val string, style *style.TextStyle) string {
	if style != nil {
		return style.Sprint(val)
	}

	return val
}

func (c *column) sprintColumnComponents(columnValueComponents []ColumnValueComponent, w int) string {
	originalLine := ""
	for _, columnValueComponent := range columnValueComponents {
		originalLine += columnValueComponent.Value
	}

	if uniseg.StringWidth(originalLine) > w {
		originalLine = removeRunes(originalLine, w)
	}

	originalLineSize := uniseg.StringWidth(originalLine)
	spaceCount := 0

	if c.leftAlign {
		originalLine = paddingRight(originalLine, w, " ")
		spaceCount = uniseg.StringWidth(originalLine) - originalLineSize
	} else {
		originalLine = fmt.Sprintf("%"+fmt.Sprintf("%v", w)+"s", originalLine)
		spaceCount = uniseg.StringWidth(originalLine) - originalLineSize
	}

	originalLineSize = uniseg.StringWidth(originalLine)

	line := ""
	lineSize := 0

	for _, columnValueComponent := range columnValueComponents {
		columnValue := columnValueComponent.Value
		lineSize += uniseg.StringWidth(columnValue)

		if lineSize <= originalLineSize {
			line += styleString(columnValue, columnValueComponent.Style)
		} else {
			// lineSize > originalLineSize
			discardSize := lineSize - originalLineSize
			columnValue = columnValue[:uniseg.StringWidth(columnValue)-discardSize]
			line += styleString(columnValue, columnValueComponent.Style)

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

func (c *column) sprintString(val string, w int) string {
	line := val

	if uniseg.StringWidth(line) > w {
		line = removeRunes(line, w)
	}

	if c.leftAlign {
		line = paddingRight(line, w, " ")
	} else {
		line = fmt.Sprintf("%"+fmt.Sprintf("%v", w)+"s", line)
	}

	return line
}

func (c *column) sprint(cv ColumnValue, w int) string {
	if w <= 0 {
		return ""
	}

	switch columnValue := cv.(type) {
	case string:
		return c.sprintString(columnValue, w)
	case []ColumnValueComponent:
		return c.sprintColumnComponents(columnValue, w)
	}

	return ""
}

// paddingRight right-pads the string with pad up to len runes
func paddingRight(str string, length int, pad string) string {
	return str + strings.Repeat(pad, length-uniseg.StringWidth(str))
}

// removeRunes truncate string length to w
func removeRunes(s string, w int) string {
	r := []rune(s)

	return string(r[:w])
}
