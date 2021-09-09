package config

import "github.com/dinhhuy258/gocui"

type Config struct {
	PathHeader  string
	TreePrefix  string
	TreeSuffix  string
	FocusPrefix string
	FocusSuffix string
	FocusBg     gocui.Attribute
	FocusFg     gocui.Attribute
}

var AppConfig *Config

func LoadConfig() {
	AppConfig = &Config{
		PathHeader:  "╭──── path",
		TreePrefix:  "├─",
		TreeSuffix:  "╰─",
		FocusPrefix: "▸[",
		FocusSuffix: "]",
		FocusBg:     gocui.ColorDefault,
		FocusFg:     gocui.ColorBlue,
	}
}
