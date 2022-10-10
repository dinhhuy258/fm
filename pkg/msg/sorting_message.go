package msg

import (
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui/key"
)

// SortByDirFirst sorts the files by directory first
func SortByDirFirst(app IApp, key key.Key, _ ...string) {
	config.AppConfig.General.Sorting.SortType = string(fs.DirFirst)

	Refresh(app, key)
}

// SortByDateModified sorts the files by date modified
func SortByDateModified(app IApp, key key.Key, _ ...string) {
	config.AppConfig.General.Sorting.SortType = string(fs.DateModified)

	Refresh(app, key)
}

// SortByName sorts the files by name
func SortByName(app IApp, key key.Key, _ ...string) {
	config.AppConfig.General.Sorting.SortType = string(fs.Name)

	Refresh(app, key)
}

// SortBySize sorts the files by size
func SortBySize(app IApp, key key.Key, _ ...string) {
	config.AppConfig.General.Sorting.SortType = string(fs.Size)

	Refresh(app, key)
}

// SortByExtension sorts the files by extension
func SortByExtension(app IApp, key key.Key, _ ...string) {
	config.AppConfig.General.Sorting.SortType = string(fs.Extension)

	Refresh(app, key)
}

// ReverseSort reverses the current sorting
func ReverseSort(app IApp, key key.Key, _ ...string) {
	*config.AppConfig.General.Sorting.Reverse = !*config.AppConfig.General.Sorting.Reverse

	Refresh(app, key)
}
