package msg

import (
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/fs"
)

// SortByDirFirst sorts the files by directory first
func SortByDirFirst(app IApp, _ ...string) {
	config.AppConfig.General.Sorting.SortType = string(fs.DirFirst)

	Refresh(app)
}

// SortByDateModified sorts the files by date modified
func SortByDateModified(app IApp, _ ...string) {
	config.AppConfig.General.Sorting.SortType = string(fs.DateModified)

	Refresh(app)
}

// SortByName sorts the files by name
func SortByName(app IApp, _ ...string) {
	config.AppConfig.General.Sorting.SortType = string(fs.Name)

	Refresh(app)
}

// SortBySize sorts the files by size
func SortBySize(app IApp, _ ...string) {
	config.AppConfig.General.Sorting.SortType = string(fs.Size)

	Refresh(app)
}

// SortByExtension sorts the files by extension
func SortByExtension(app IApp, _ ...string) {
	config.AppConfig.General.Sorting.SortType = string(fs.Extension)

	Refresh(app)
}

// ReverseSort reverses the current sorting
func ReverseSort(app IApp, _ ...string) {
	config.AppConfig.General.Sorting.Reverse = !config.AppConfig.General.Sorting.Reverse

	Refresh(app)
}
