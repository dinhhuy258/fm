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
	explorerTableConfig := config.AppConfig.General.ExplorerTable

	indexHeaderTextStyle := style.FromStyleConfig(explorerTableConfig.IndexHeader.Style)
	nameHeaderTextStyle := style.FromStyleConfig(explorerTableConfig.NameHeader.Style)
	permissionsHeaderTextStyle := style.FromStyleConfig(explorerTableConfig.PermissionsHeader.Style)
	sizeHeaderTextStyle := style.FromStyleConfig(explorerTableConfig.SizeHeader.Style)

	x, _ := ehv.Size()
	ehv.headerRow.SetWidth(x)

	rowString, err := ehv.headerRow.Sprint(
		[]ColumnValue{
			[]ColumnValueComponent{
				{
					Value: explorerTableConfig.IndexHeader.Name,
					Style: &indexHeaderTextStyle,
				},
			},
			[]ColumnValueComponent{
				{
					Value: explorerTableConfig.NameHeader.Name,
					Style: &nameHeaderTextStyle,
				},
			},
			[]ColumnValueComponent{
				{
					Value: explorerTableConfig.PermissionsHeader.Name,
					Style: &permissionsHeaderTextStyle,
				},
			},
			[]ColumnValueComponent{
				{
					Value: explorerTableConfig.SizeHeader.Name,
					Style: &sizeHeaderTextStyle,
				},
			},
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

	explorerRow *Row
	icons       nodeTypes

	defaultFileTextStyle        style.TextStyle
	focusFileTextStyle          style.TextStyle
	selectionFileTextStyle      style.TextStyle
	focusSelectionFileTextStyle style.TextStyle

	defaultDirectoryTextStyle        style.TextStyle
	focusDirectoryTextStyle          style.TextStyle
	selectionDirectoryTextStyle      style.TextStyle
	focusSelectionDirectoryTextStyle style.TextStyle
}

func newExplorerView(v *gocui.View) *ExplorerView {
	cfg := config.AppConfig

	ev := &ExplorerView{
		View:        newView(v),
		explorerRow: newExplorerRow(),
	}

	ev.Frame = false

	ev.defaultFileTextStyle = style.FromStyleConfig(cfg.General.DefaultUI.FileStyle)
	ev.focusFileTextStyle = style.FromStyleConfig(cfg.General.FocusUI.FileStyle)
	ev.focusSelectionFileTextStyle = style.FromStyleConfig(cfg.General.FocusSelectionUI.FileStyle)
	ev.selectionFileTextStyle = style.FromStyleConfig(cfg.General.SelectionUI.FileStyle)

	ev.defaultDirectoryTextStyle = style.FromStyleConfig(cfg.General.DefaultUI.DirectoryStyle)
	ev.focusDirectoryTextStyle = style.FromStyleConfig(cfg.General.FocusUI.DirectoryStyle)
	ev.focusSelectionDirectoryTextStyle = style.FromStyleConfig(cfg.General.FocusSelectionUI.DirectoryStyle)
	ev.selectionDirectoryTextStyle = style.FromStyleConfig(cfg.General.SelectionUI.DirectoryStyle)

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
	explorerTableConfig := config.AppConfig.General.ExplorerTable

	r := &Row{}

	r.AddColumn(explorerTableConfig.IndexHeader.Percentage, true)
	r.AddColumn(explorerTableConfig.NameHeader.Percentage, true)
	r.AddColumn(explorerTableConfig.PermissionsHeader.Percentage, true)
	r.AddColumn(explorerTableConfig.SizeHeader.Percentage, false)

	return r
}

func (ev *ExplorerView) UpdateView(entries []fs.IEntry, selections set.Set[string], focus int) {
	entriesSize := len(entries)
	lines := make([]string, entriesSize)
	cfg := config.AppConfig

	for idx, entry := range entries {
		isEntryFocused := idx == focus
		isEntrySelected := selections.Contains(entry.GetPath())

		entryIcon := ev.getEntryIcon(entry, isEntryFocused, isEntrySelected)

		name := entry.GetName()

		var prefix, suffix, entryTreePrefix string

		var entryTextStyle style.TextStyle

		switch {
		case isEntryFocused && isEntrySelected:
			prefix = cfg.General.FocusSelectionUI.Prefix
			suffix = cfg.General.FocusSelectionUI.Suffix

			if entry.IsDirectory() {
				entryTextStyle = ev.focusSelectionDirectoryTextStyle
			} else {
				entryTextStyle = ev.focusSelectionFileTextStyle
			}
		case isEntryFocused:
			prefix = cfg.General.FocusUI.Prefix
			suffix = cfg.General.FocusUI.Suffix

			if entry.IsDirectory() {
				entryTextStyle = ev.focusDirectoryTextStyle
			} else {
				entryTextStyle = ev.focusFileTextStyle
			}
		case isEntrySelected:
			prefix = cfg.General.SelectionUI.Prefix
			suffix = cfg.General.SelectionUI.Suffix

			if entry.IsDirectory() {
				entryTextStyle = ev.selectionDirectoryTextStyle
			} else {
				entryTextStyle = ev.selectionFileTextStyle
			}
		default:
			prefix = cfg.General.DefaultUI.Prefix
			suffix = cfg.General.DefaultUI.Suffix

			if entry.IsDirectory() {
				entryTextStyle = ev.defaultDirectoryTextStyle
			} else {
				entryTextStyle = ev.defaultFileTextStyle
			}
		}

		switch {
		case idx == entriesSize-1:
			entryTreePrefix = cfg.General.ExplorerTable.LastEntryPrefix
		case idx == 0:
			entryTreePrefix = cfg.General.ExplorerTable.FirstEntryPrefix
		default:
			entryTreePrefix = cfg.General.ExplorerTable.EntryPrefix
		}

		index := strconv.Itoa(idx + 1)
		fileMode := entry.GetFileMode()
		size := fs.Humanize(entry.GetSize())

		line, err := ev.explorerRow.Sprint([]ColumnValue{
			index,
			[]ColumnValueComponent{
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
			},
			fileMode,
			size,
		})
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

func (ev *ExplorerView) getEntryIcon(entry fs.IEntry, isEntryFocused, isEntrySelected bool) nodeType {
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

	switch {
	case isEntrySelected && isEntryFocused:
		icon.style = ev.focusDirectoryTextStyle
	case isEntrySelected:
		icon.style = ev.selectionFileTextStyle
	case isEntryFocused:
		icon.style = ev.focusDirectoryTextStyle
	}

	return icon
}
