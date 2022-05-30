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
		[]string{
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

type icon struct {
	val   string
	style style.TextStyle
}

func (i icon) sprint() string {
	return i.style.Sprint(i.val)
}

type icons struct {
	file       icon
	directory  icon
	extensions map[string]icon
}

type ExplorerView struct {
	*View

	explorerRow        *style.Row
	icons              icons
	fileTextStyle      style.TextStyle
	directoryTextStyle style.TextStyle
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

	ev.icons = icons{
		file: icon{
			val:   cfg.NodeTypesConfig.File.Icon,
			style: ev.fileTextStyle,
		},
		directory: icon{
			val:   cfg.NodeTypesConfig.Directory.Icon,
			style: ev.directoryTextStyle,
		},
		extensions: map[string]icon{},
	}

	for ext, ntc := range cfg.NodeTypesConfig.Extensions {
		if ntc.Color != "" {
			ev.icons.extensions[ext] = icon{
				val:   ntc.Icon,
				style: style.FromBasicFg(style.StringToColor(ntc.Color)),
			}
		} else {
			ev.icons.extensions[ext] = icon{
				val:   ntc.Icon,
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
		var nameTextStyle *style.TextStyle
		var icon icon

		if entry.IsDirectory() {
			nameTextStyle = &ev.directoryTextStyle
			icon = ev.icons.directory
		} else {
			nameTextStyle = &ev.fileTextStyle
			icon = ev.icons.file
		}

		if i, hasIcon := ev.icons.extensions[entry.GetExt()]; hasIcon {
			icon = i
		}

		name := entry.GetName()
		var prefix string
		var suffix string
		var rowPrefix string

		switch {
		case idx == focus:
			prefix = cfg.FocusPrefix
			suffix = cfg.FocusSuffix
		case selections.Contains(entry.GetPath()):
			prefix = cfg.SelectionPrefix
			suffix = cfg.SelectionSuffix
		default:
			// TODO: Configurate these values
			prefix = "  "
			suffix = ""
		}

		if idx == entriesSize-1 {
			rowPrefix = cfg.PathSuffix
		} else {
			rowPrefix = cfg.PathPrefix
		}

		index := strconv.Itoa(idx + 1)
		fileMode := entry.GetFileMode()
		size := fs.Humanize(entry.GetSize())

		path := rowPrefix + icon.sprint() + nameTextStyle.Sprint(prefix+name+suffix)
		line, err := ev.explorerRow.Sprint([]string{index, path, fileMode, size})
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
