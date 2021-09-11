package config

import (
	"github.com/dinhhuy258/fm/pkg/style"
	"github.com/dinhhuy258/gocui"
	"github.com/gookit/color"
)

type Config struct {
	IndexHeader     string
	IndexPercentage int
	PathHeader      string
	PathPercentage  int
	SizeHeader      string
	SizePercentage  int
	PathPrefix      string
	PathSuffix      string
	FocusPrefix     string
	FocusSuffix     string
	FocusBg         gocui.Attribute
	FocusFg         gocui.Attribute
	SelectionPrefix string
	SelectionSuffix string
	SelectionStyle  style.TextStyle
	FolderIcon      string
	FileIcon        string
	DirectoryStyle  style.TextStyle
	SizeStyle       style.TextStyle
}

var AppConfig *Config

func LoadConfig() {
	AppConfig = &Config{
		IndexHeader:     "index",
		IndexPercentage: 10,
		PathHeader:      "╭──── path",
		PathPercentage:  70,
		SizeHeader:      "size",
		SizePercentage:  20,
		PathPrefix:      "├─",
		PathSuffix:      "╰─",
		FocusPrefix:     "▸[",
		FocusSuffix:     "]",
		FocusBg:         gocui.ColorDefault,
		FocusFg:         gocui.ColorBlue,
		SelectionPrefix: "{",
		SelectionSuffix: "}",
		SelectionStyle:  style.FromBasicFg(color.Green),
		FolderIcon:      "",
		FileIcon:        "",
		DirectoryStyle:  style.FromBasicFg(color.Cyan),
		SizeStyle:       style.FromBasicFg(color.White),
	}
}
