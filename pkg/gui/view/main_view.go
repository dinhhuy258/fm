package view

import (
	"log"
	"strconv"

	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui/view/row"
	"github.com/dinhhuy258/fm/pkg/gui/view/style"
	"github.com/dinhhuy258/gocui"
	"github.com/gookit/color"
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
		directoryRow: newRow(&config.AppConfig.DirectoryColor),
		selectionRow: newRow(&config.AppConfig.SelectionColor),
	}

	mv.v.v.Frame = false
	mv.v.v.Highlight = true
	mv.v.v.SelBgColor = config.AppConfig.FocusBg
	mv.v.v.SelFgColor = config.AppConfig.FocusFg

	return mv
}

func newRow(pathColor *color.Color) *row.Row {
	r := &row.Row{}
	r.AddCell(config.AppConfig.IndexPercentage, true, nil)

	if pathColor != nil {
		pathStyle := style.FromBasicFg(*pathColor)
		r.AddCell(config.AppConfig.PathPercentage, true, &pathStyle)
	} else {
		r.AddCell(config.AppConfig.PathPercentage, true, nil)
	}

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

func (mv *MainView) RenderDir(dir *fs.Directory, selections map[string]struct{}, focusIdx int) {
	visibleNodeSize := len(dir.VisibleNodes)
	lines := make([]string, visibleNodeSize)
	cfg := config.AppConfig

	for i, node := range dir.VisibleNodes {
		fileIcon := cfg.FileIcon + " "
		if node.IsDir {
			fileIcon = cfg.FolderIcon + " "
		}

		_, isSelected := selections[node.AbsolutePath]

		var path string

		switch {
		case i == focusIdx:
			path = cfg.FocusPrefix + fileIcon + node.RelativePath + cfg.FocusSuffix
		case isSelected:
			path = cfg.SelectionPrefix + fileIcon + node.RelativePath + cfg.SelectionSuffix
		default:
			path = "  " + fileIcon + node.RelativePath
		}

		if i == visibleNodeSize-1 {
			path = cfg.PathSuffix + path
		} else {
			path = cfg.PathPrefix + path
		}

		r := mv.fileRow
		if isSelected {
			r = mv.selectionRow
		} else if node.IsDir {
			r = mv.directoryRow
		}

		size := fs.Humanize(node.Size)
		index := strconv.Itoa(i-focusIdx) + "|" + strconv.Itoa(i)

		line, err := r.Sprint([]string{index, path, size})
		if err != nil {
			log.Fatalf("failed to sprint row %v", err)
		}

		lines[i] = line
	}

	mv.v.SetViewContent(lines)
}

func (mv *MainView) SetTitle(path string, numberOfFiles int) {
	mv.hv.v.Title = (" " + path + " (" + strconv.Itoa(numberOfFiles) + ") ")
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

func (mv *MainView) SetAsCurrentView() {
	_, err := mv.v.g.SetCurrentView(mv.v.v.Name())
	if err != nil {
		log.Fatalf("failed to set main view as the current view %v", err)
	}
}
