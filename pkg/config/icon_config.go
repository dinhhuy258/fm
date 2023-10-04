package config

// getSpecialsNodeTypeConfig get default configuration for the special node type.
func getSpecialsNodeTypeConfig() map[string]*NodeTypeConfig {
	specialsNodeTypeConfig := make(map[string]*NodeTypeConfig)
	for fileName, icon := range IconsByFilename {
		specialsNodeTypeConfig[fileName] = &NodeTypeConfig{
			Icon: icon.Glyph,
			Style: &StyleConfig{
				Fg: icon.Color,
			},
		}
	}

	for dirName, icon := range IconsByDir {
		specialsNodeTypeConfig[dirName] = &NodeTypeConfig{
			Icon: icon.Glyph,
			Style: &StyleConfig{
				Fg: icon.Color,
			},
		}
	}

	return specialsNodeTypeConfig
}

// getExtensionsNodeTypeConfig get default configuration for the extensions node type.
func getExtensionsNodeTypeConfig() map[string]*NodeTypeConfig {
	extensionsNodeTypeConfig := make(map[string]*NodeTypeConfig)
	for ext, icon := range IconsByExtension {
		extensionsNodeTypeConfig[ext] = &NodeTypeConfig{
			Icon: icon.Glyph,
			Style: &StyleConfig{
				Fg: icon.Color,
			},
		}
	}

	return extensionsNodeTypeConfig
}
