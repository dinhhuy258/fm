package row

import (
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/style"
)

type MainRow struct {
	HeaderRow    *Row
	FileRow      *Row
	DirectoryRow *Row
	SelectionRow *Row
}

func (mr *MainRow) SetWidth(width int) {
	mr.HeaderRow.SetWidth(width)
	mr.FileRow.SetWidth(width)
	mr.DirectoryRow.SetWidth(width)
	mr.SelectionRow.SetWidth(width)
}

func NewMainRow() *MainRow {
	return &MainRow{
		HeaderRow:    newRow(nil),
		FileRow:      newRow(nil),
		DirectoryRow: newRow(&config.AppConfig.DirectoryStyle),
		SelectionRow: newRow(&config.AppConfig.SelectionStyle),
	}
}

func newRow(pathStyle *style.TextStyle) *Row {
	r := &Row{}
	r.AddCell(config.AppConfig.IndexPercentage, true, nil)
	r.AddCell(config.AppConfig.PathPercentage, true, pathStyle)
	r.AddCell(config.AppConfig.SizePercentage, false, nil)

	return r
}
