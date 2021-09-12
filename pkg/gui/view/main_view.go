package view

import (
	"strconv"

	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/row"
	"github.com/dinhhuy258/fm/pkg/style"
	"github.com/dinhhuy258/gocui"
)

type MainView struct {
	v            *View
	hv           *View
	headerRow    *row.Row
	fileRow      *row.Row
	directoryRow *row.Row
	selectionRow *row.Row
}

func newMainView(g *gocui.Gui, v *gocui.View, hv *gocui.View) *MainView {
	mv := &MainView{
		v:            newView(g, v),
		hv:           newView(g, hv),
		headerRow:    newRow(nil),
		fileRow:      newRow(nil),
		directoryRow: newRow(&config.AppConfig.DirectoryStyle),
		selectionRow: newRow(&config.AppConfig.SelectionStyle),
	}

	mv.v.v.Frame = false
	mv.v.v.Highlight = true
	mv.v.v.SelBgColor = config.AppConfig.FocusBg
	mv.v.v.SelFgColor = config.AppConfig.FocusFg

	return mv
}

func newRow(pathStyle *style.TextStyle) *row.Row {
	r := &row.Row{}
	r.AddCell(config.AppConfig.IndexPercentage, true, nil)
	r.AddCell(config.AppConfig.PathPercentage, true, pathStyle)
	r.AddCell(config.AppConfig.SizePercentage, false, nil)

	return r
}

func (mv *MainView) layout() error {
	x, _ := mv.v.v.Size()
	mv.headerRow.SetWidth(x)
	mv.directoryRow.SetWidth(x)
	mv.fileRow.SetWidth(x)
	mv.selectionRow.SetWidth(x)

	rowString, err := mv.headerRow.Sprint(
		[]string{config.AppConfig.IndexHeader, config.AppConfig.PathHeader, config.AppConfig.SizeHeader},
	)
	if err != nil {
		return err
	}

	mv.hv.SetViewContent([]string{rowString})

	return nil
}

func (mv *MainView) RenderDir(dir *fs.Directory, selections map[string]struct{}, focusIdx int) error {
	nodeSize := len(dir.Nodes)
	lines := make([]string, nodeSize)
	config := config.AppConfig

	for i, node := range dir.Nodes {
		fileIcon := config.FileIcon + " "
		if node.IsDir {
			fileIcon = config.FolderIcon + " "
		}

		_, isSelected := selections[node.AbsolutePath]

		var path string

		switch {
		case i == focusIdx:
			path = config.FocusPrefix + fileIcon + node.RelativePath + config.FocusSuffix
		case isSelected:
			path = config.SelectionPrefix + fileIcon + node.RelativePath + config.SelectionSuffix
		default:
			path = "  " + fileIcon + node.RelativePath
		}

		if i == nodeSize-1 {
			path = config.PathSuffix + path
		} else {
			path = config.PathPrefix + path
		}

		row := mv.fileRow
		if isSelected {
			row = mv.selectionRow
		} else if node.IsDir {
			row = mv.directoryRow
		}

		size := fs.Humanize(node.Size)
		index := strconv.Itoa(i-focusIdx) + "|" + strconv.Itoa(i)

		line, err := row.Sprint([]string{index, path, size})
		if err != nil {
			return err
		}

		lines[i] = line
	}

	mv.v.SetViewContent(lines)

	return nil
}

func (mv *MainView) SetTitle(header string) {
	mv.hv.v.Title = header
}

func (mv *MainView) SetOrigin(x, y int) error {
	return mv.v.SetOrigin(x, y)
}

func (mv *MainView) SetCursor(x, y int) error {
	return mv.v.SetCursor(x, y)
}

func (mv *MainView) NextCursor() error {
	return mv.v.NextCursor()
}

func (mv *MainView) PreviousCursor() error {
	return mv.v.PreviousCursor()
}

func (mv *MainView) SetAsCurrentView() error {
	_, err := mv.v.g.SetCurrentView(mv.v.v.Name())

	return err
}
