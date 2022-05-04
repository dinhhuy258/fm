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

func (mv *MainView) RenderEntries(entries []fs.IEntry, selections map[string]struct{}, focus int) {
	entriesSize := len(entries)
	lines := make([]string, entriesSize)
	cfg := config.AppConfig

	for idx, entry := range entries {
		fileIcon := cfg.FileIcon + " "
		if entry.IsDirectory() {
			fileIcon = cfg.FolderIcon + " "
		}

		_, isSelected := selections[entry.GetPath()]

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

	mv.v.SetViewContent(lines)
}

func (mv *MainView) RenderDir(nodes []*fs.Node, selections map[string]struct{}, focus int) {
	nodesSize := len(nodes)
	lines := make([]string, nodesSize)
	cfg := config.AppConfig

	for i, node := range nodes {
		fileIcon := cfg.FileIcon + " "
		if node.IsDir {
			fileIcon = cfg.FolderIcon + " "
		}

		_, isSelected := selections[node.AbsolutePath]

		var path string

		switch {
		case i == focus:
			path = cfg.FocusPrefix + fileIcon + node.RelativePath + cfg.FocusSuffix
		case isSelected:
			path = cfg.SelectionPrefix + fileIcon + node.RelativePath + cfg.SelectionSuffix
		default:
			path = "  " + fileIcon + node.RelativePath
		}

		if i == nodesSize-1 {
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
		index := strconv.Itoa(i-focus) + "|" + strconv.Itoa(i)

		line, err := r.Sprint([]string{index, path, size})
		if err != nil {
			log.Fatalf("failed to sprint row %v", err)
		}

		lines[i] = line
	}

	mv.v.SetViewContent(lines)
}

func (mv *MainView) SetTitle(title string) {
	mv.hv.v.Title = title
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
