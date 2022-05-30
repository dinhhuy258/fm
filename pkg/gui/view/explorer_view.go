package view

import (
	"log"
	"strconv"
	"strings"

	set "github.com/deckarep/golang-set/v2"
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui/view/style"
	"github.com/dinhhuy258/gocui"
)

type ExplorerHeaderView struct {
	*View

	headerRow *style.Row
}

func newExplorerHeaderView(v *gocui.View) *ExplorerHeaderView {
	ehv := &ExplorerHeaderView{
		View:      newView(v),
		headerRow: newExplorerRow(),
	}

	return ehv
}

func (ehv *ExplorerHeaderView) layout() error {
	x, _ := ehv.Size()
	ehv.headerRow.SetWidth(x)

	rowString, err := ehv.headerRow.Sprint(
		[]style.CellValue{
			config.AppConfig.IndexHeader,
			config.AppConfig.PathHeader,
			config.AppConfig.FileModeHeader,
			config.AppConfig.SizeHeader,
		},
	)
	if err != nil {
		return err
	}

	ehv.SetContent(rowString)

	return nil
}

type nodeType struct {
	icon  string
	style style.TextStyle
}

type nodeTypes struct {
	file       nodeType
	directory  nodeType
	extensions map[string]nodeType
}

type ExplorerView struct {
	*View

	explorerRow        *style.Row
	icons              nodeTypes
	fileTextStyle      style.TextStyle
	directoryTextStyle style.TextStyle
	selectionTextStyle style.TextStyle
}

func newExplorerView(v *gocui.View) *ExplorerView {
	cfg := config.AppConfig

	ev := &ExplorerView{
		View:        newView(v),
		explorerRow: newExplorerRow(),
	}

	ev.Frame = false
	ev.Highlight = true

	ev.SelBgColor = style.StringToGoCuiColor(config.AppConfig.FocusBg)
	ev.SelFgColor = style.StringToGoCuiColor(config.AppConfig.FocusFg)

	ev.directoryTextStyle = style.FromBasicFg(style.StringToColor(cfg.NodeTypesConfig.Directory.Color))
	ev.fileTextStyle = style.FromBasicFg(style.StringToColor(cfg.NodeTypesConfig.File.Color))
	ev.selectionTextStyle = style.FromBasicFg(style.StringToColor(cfg.SelectionColor))

	ev.icons = nodeTypes{
		file: nodeType{
			icon:  cfg.NodeTypesConfig.File.Icon,
			style: ev.fileTextStyle,
		},
		directory: nodeType{
			icon:  cfg.NodeTypesConfig.Directory.Icon,
			style: ev.directoryTextStyle,
		},
		extensions: map[string]nodeType{},
	}

	for ext, ntc := range cfg.NodeTypesConfig.Extensions {
		if ntc.Color != "" {
			ev.icons.extensions[ext] = nodeType{
				icon:  ntc.Icon,
				style: style.FromBasicFg(style.StringToColor(ntc.Color)),
			}
		} else {
			ev.icons.extensions[ext] = nodeType{
				icon:  ntc.Icon,
				style: ev.fileTextStyle,
			}
		}
	}

	return ev
}

func newExplorerRow() *style.Row {
	r := &style.Row{}

	r.AddCell(config.AppConfig.IndexPercentage, true)
	r.AddCell(config.AppConfig.PathPercentage, true)
	r.AddCell(config.AppConfig.FileModePercentage, true)
	r.AddCell(config.AppConfig.SizePercentage, false)

	return r
}

func (ev *ExplorerView) UpdateView(entries []fs.IEntry, selections set.Set[string], focus int) {
	entriesSize := len(entries)
	lines := make([]string, entriesSize)
	cfg := config.AppConfig

	for idx, entry := range entries {
		isEntrySelected := selections.Contains(entry.GetPath())

		entryTextStyle := ev.getEntryTextStyle(entry, isEntrySelected)
		entryIcon := ev.getEntryIcon(entry, isEntrySelected)

		name := entry.GetName()

		var prefix, suffix, entryTreePrefix string

		switch {
		case idx == focus:
			prefix = cfg.FocusPrefix
			suffix = cfg.FocusSuffix
		case isEntrySelected:
			prefix = cfg.SelectionPrefix
			suffix = cfg.SelectionSuffix
		default:
			// TODO: Configurate these values
			prefix = "  "
			suffix = ""
		}

		if idx == entriesSize-1 {
			entryTreePrefix = cfg.PathSuffix
		} else {
			entryTreePrefix = cfg.PathPrefix
		}

		index := strconv.Itoa(idx + 1)
		fileMode := entry.GetFileMode()
		size := fs.Humanize(entry.GetSize())

		line, err := ev.explorerRow.Sprint([]style.CellValue{index, []style.CellValueComponent{
			{
				Value: entryTreePrefix,
				Style: nil,
			},
			{
				Value: prefix,
				Style: &entryTextStyle,
			},
			{
				Value: entryIcon.icon,
				Style: &entryIcon.style,
			},
			{
				Value: " ",
				Style: nil,
			},
			{
				Value: name,
				Style: &entryTextStyle,
			},
			{
				Value: suffix,
				Style: &entryTextStyle,
			},
		}, fileMode, size})
		if err != nil {
			log.Fatalf("failed to sprint row %v", err)
		}

		lines[idx] = line
	}

	ev.SetContent(strings.Join(lines, "\n"))
}

func (ev *ExplorerView) layout() error {
	x, _ := ev.Size()
	ev.explorerRow.SetWidth(x)

	return nil
}

func (ev *ExplorerView) getEntryTextStyle(entry fs.IEntry, isEntrySelected bool) style.TextStyle {
	switch {
	case isEntrySelected:
		return ev.selectionTextStyle
	case entry.IsDirectory():
		return ev.directoryTextStyle
	default:
		return ev.fileTextStyle
	}
}

func (ev *ExplorerView) getEntryIcon(entry fs.IEntry, isEntrySelected bool) nodeType {
	var icon nodeType

	if i, hasIcon := ev.icons.extensions[entry.GetExt()]; hasIcon {
		icon = i
	} else if entry.IsDirectory() {
		icon = ev.icons.directory
	} else {
		icon = ev.icons.file
	}

	if isEntrySelected {
		icon.style = ev.selectionTextStyle
	}

	return icon
}
