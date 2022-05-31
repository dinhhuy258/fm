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

	headerRow *Row
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
		[]ColumnValue{
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

	explorerRow             *Row
	icons                   nodeTypes
	defaultTextStyle        style.TextStyle
	selectionTextStyle      style.TextStyle
	focusTextStyle          style.TextStyle
	focusSelectionTextStyle style.TextStyle
}

func newExplorerView(v *gocui.View) *ExplorerView {
	cfg := config.AppConfig

	ev := &ExplorerView{
		View:        newView(v),
		explorerRow: newExplorerRow(),
	}

	ev.Frame = false
	// ev.Highlight = true
	// ev.SelBgColor = gocui.GetColor(config.AppConfig.General.SelectionUI.Style.Bg)
	// ev.SelFgColor = gocui.GetColor(config.AppConfig.General.SelectionUI.Style.Fg)

	ev.defaultTextStyle = style.FromStyleConfig(cfg.General.DefaultUI.Style)
	ev.focusTextStyle = style.FromStyleConfig(cfg.General.FocusUI.Style)
	ev.focusSelectionTextStyle = style.FromStyleConfig(cfg.General.FocusSelectionUI.Style)
	ev.selectionTextStyle = style.FromStyleConfig(cfg.General.SelectionUI.Style)

	ev.icons = nodeTypes{
		file: nodeType{
			icon:  cfg.NodeTypesConfig.File.Icon,
			style: style.FromStyleConfig(cfg.NodeTypesConfig.File.Style),
		},
		directory: nodeType{
			icon:  cfg.NodeTypesConfig.Directory.Icon,
			style: style.FromStyleConfig(cfg.NodeTypesConfig.Directory.Style),
		},
		extensions: map[string]nodeType{},
	}

	for ext, ntc := range cfg.NodeTypesConfig.Extensions {
		ev.icons.extensions[ext] = nodeType{
			icon:  ntc.Icon,
			style: style.FromStyleConfig(ntc.Style),
		}
	}

	return ev
}

func newExplorerRow() *Row {
	r := &Row{}

	r.AddColumn(config.AppConfig.IndexPercentage, true)
	r.AddColumn(config.AppConfig.PathPercentage, true)
	r.AddColumn(config.AppConfig.FileModePercentage, true)
	r.AddColumn(config.AppConfig.SizePercentage, false)

	return r
}

func (ev *ExplorerView) UpdateView(entries []fs.IEntry, selections set.Set[string], focus int) {
	entriesSize := len(entries)
	lines := make([]string, entriesSize)
	cfg := config.AppConfig

	for idx, entry := range entries {
		isEntryFocused := idx == focus
		isEntrySelected := selections.Contains(entry.GetPath())

		entryIcon := ev.getEntryIcon(entry, isEntrySelected)

		name := entry.GetName()

		var prefix, suffix, entryTreePrefix string

		var entryTextStyle style.TextStyle

		switch {
		case isEntryFocused && isEntrySelected:
			prefix = cfg.General.FocusSelectionUI.Prefix
			suffix = cfg.General.FocusSelectionUI.Suffix

			entryTextStyle = ev.focusSelectionTextStyle
		case isEntryFocused:
			prefix = cfg.General.FocusUI.Prefix
			suffix = cfg.General.FocusUI.Suffix

			entryTextStyle = ev.focusTextStyle
		case isEntrySelected:
			prefix = cfg.General.SelectionUI.Prefix
			suffix = cfg.General.SelectionUI.Suffix

			entryTextStyle = ev.selectionTextStyle
		default:
			prefix = cfg.General.DefaultUI.Prefix
			suffix = cfg.General.DefaultUI.Suffix

			entryTextStyle = ev.defaultTextStyle
		}

		if idx == entriesSize-1 {
			entryTreePrefix = cfg.PathSuffix
		} else {
			entryTreePrefix = cfg.PathPrefix
		}

		index := strconv.Itoa(idx + 1)
		fileMode := entry.GetFileMode()
		size := fs.Humanize(entry.GetSize())

		line, err := ev.explorerRow.Sprint([]ColumnValue{index, []ColumnValueComponent{
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

func (ev *ExplorerView) getEntryIcon(entry fs.IEntry, isEntrySelected bool) nodeType {
	var icon nodeType

	extensionIcon, hasExtIcon := ev.icons.extensions[entry.GetExt()]

	switch {
	case hasExtIcon:
		icon = extensionIcon
	case entry.IsDirectory():
		icon = ev.icons.directory
	default:
		icon = ev.icons.file
	}

	if isEntrySelected {
		icon.style = ev.selectionTextStyle
	}

	return icon
}
