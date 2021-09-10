package row

import (
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/style"
)

type MainRow struct {
	HeaderRow    *Row
	FileRow      *Row
	DirectoryRow *Row
}

func NewMainRow() *MainRow {
	return &MainRow{
		HeaderRow:    newRow(nil, nil),
		FileRow:      newRow(nil, nil),
		DirectoryRow: newRow(&config.AppConfig.DirectoryStyle, nil),
	}
}

func newRow(pathStyle *style.TextStyle, sizeStyle *style.TextStyle) *Row {
	r := &Row{}
	r.AddCell(config.AppConfig.PathPercentage, true, pathStyle)
	r.AddCell(config.AppConfig.SizePercentage, false, sizeStyle)

	return r
}
