package view

import (
	"log"
	"strconv"

	set "github.com/deckarep/golang-set/v2"
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui/view/style"
	"github.com/dinhhuy258/fm/pkg/optional"
	"github.com/dinhhuy258/gocui"
	"github.com/gookit/color"
)

type ExplorerHeaderView struct {
	*View
	headerRow *style.Row
}

func newExplorerHeaderView(v *gocui.View) *ExplorerHeaderView {
	ehv := &ExplorerHeaderView{
		View:      newView(v),
		headerRow: newRow(optional.NewEmpty[color.Color]()),
	}

	return ehv
}

func (ehv *ExplorerHeaderView) layout() error {
	x, _ := ehv.Size()
	ehv.headerRow.SetWidth(x)

	rowString, err := ehv.headerRow.Sprint(
		[]string{config.AppConfig.IndexHeader, config.AppConfig.PathHeader, config.AppConfig.SizeHeader},
	)
	if err != nil {
		return err
	}

	ehv.SetViewContent([]string{rowString})

	return nil
}

type ExplorerView struct {
	*View
	fileRow      *style.Row
	directoryRow *style.Row
	selectionRow *style.Row
}

func newExplorerView(v *gocui.View) *ExplorerView {
	ev := &ExplorerView{
		View:         newView(v),
		fileRow:      newRow(optional.NewEmpty[color.Color]()),
		directoryRow: newRow(optional.New(style.StringToColor(config.AppConfig.DirectoryColor))),
		selectionRow: newRow(optional.New(style.StringToColor(config.AppConfig.SelectionColor))),
	}

	ev.Frame = false
	ev.Highlight = true
	ev.SelBgColor = style.StringToGoCuiColor(config.AppConfig.FocusBg)
	ev.SelFgColor = style.StringToGoCuiColor(config.AppConfig.FocusFg)

	return ev
}

func newRow(pathColor optional.Optional[color.Color]) *style.Row {
	r := &style.Row{}
	r.AddCell(config.AppConfig.IndexPercentage, true, nil)

	pathColor.IfPresentOrElse(func(c *color.Color) {
		pathStyle := style.FromBasicFg(*c)
		r.AddCell(config.AppConfig.PathPercentage, true, &pathStyle)
	}, func() {
		r.AddCell(config.AppConfig.PathPercentage, true, nil)
	})

	r.AddCell(config.AppConfig.SizePercentage, false, nil)

	return r
}

func (ev *ExplorerView) UpdateView(entries []fs.IEntry, selections set.Set[string], focus int) {
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

		r := ev.fileRow
		if isSelected {
			r = ev.selectionRow
		} else if entry.IsDirectory() {
			r = ev.directoryRow
		}

		size := fs.Humanize(entry.GetSize())
		index := strconv.Itoa(idx-focus) + "|" + strconv.Itoa(idx)

		line, err := r.Sprint([]string{index, path, size})
		if err != nil {
			log.Fatalf("failed to sprint row %v", err)
		}

		lines[idx] = line
	}

	ev.SetViewContent(lines)
}

func (ev *ExplorerView) layout() error {
	x, _ := ev.Size()
	ev.directoryRow.SetWidth(x)
	ev.fileRow.SetWidth(x)
	ev.selectionRow.SetWidth(x)

	return nil
}
