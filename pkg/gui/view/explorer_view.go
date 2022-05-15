package view

import (
	"log"
	"strconv"

	set "github.com/deckarep/golang-set/v2"
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui/view/row"
	"github.com/dinhhuy258/fm/pkg/gui/view/style"
	"github.com/dinhhuy258/gocui"
	"github.com/gookit/color"
)

type ExplorerView struct {
	*View
	hv           *View
	headerRow    *row.Row
	fileRow      *row.Row
	directoryRow *row.Row
	selectionRow *row.Row
}

func newExplorerView(g *gocui.Gui, v *gocui.View, hv *gocui.View) *ExplorerView {
	mv := &ExplorerView{
		View:         newView(g, v),
		hv:           newView(g, hv),
		headerRow:    newRow(nil),
		fileRow:      newRow(nil),
		directoryRow: newRow(&config.AppConfig.DirectoryColor),
		selectionRow: newRow(&config.AppConfig.SelectionColor),
	}

	mv.v.Frame = false
	mv.v.Highlight = true
	mv.v.SelBgColor = config.AppConfig.FocusBg
	mv.v.SelFgColor = config.AppConfig.FocusFg

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

func (mv *ExplorerView) UpdateView(entries []fs.IEntry, selections set.Set[string], focus int) {
	entriesSize := len(entries)
	lines := make([]string, entriesSize)
	cfg := config.AppConfig

	for idx, entry := range entries {
		fileIcon := cfg.FileIcon + " "
		if entry.IsDirectory() {
			fileIcon = cfg.FolderIcon + " "
		}

		isSelected := selections.Contains(entry.GetPath())

		var path string

		switch {
		case idx == focus:
			path = cfg.FocusPrefix + fileIcon + entry.GetName() + cfg.FocusSuffix
		case isSelected:
			path = cfg.SelectionPrefix + fileIcon + entry.GetName() + cfg.SelectionSuffix
		default:
			path = "  " + fileIcon + entry.GetName()
		}

		if idx == entriesSize-1 {
			path = cfg.PathSuffix + path
		} else {
			path = cfg.PathPrefix + path
		}

		r := mv.fileRow
		if isSelected {
			r = mv.selectionRow
		} else if entry.IsDirectory() {
			r = mv.directoryRow
		}

		size := fs.Humanize(entry.GetSize())
		index := strconv.Itoa(idx-focus) + "|" + strconv.Itoa(idx)

		line, err := r.Sprint([]string{index, path, size})
		if err != nil {
			log.Fatalf("failed to sprint row %v", err)
		}

		lines[idx] = line
	}

	mv.SetViewContent(lines)
}

func (mv *ExplorerView) SetTitle(title string) {
	mv.hv.v.Title = title
}

func (mv *ExplorerView) SetAsCurrentView() {
	_, err := mv.g.SetCurrentView(mv.v.Name())
	if err != nil {
		log.Fatalf("failed to set explorer view as the current view %v", err)
	}
}

func (mv *ExplorerView) layout() error {
	x, _ := mv.v.Size()
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
