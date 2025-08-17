package types

// SortType represents different sorting algorithms for entries
type SortType string

const (
	// SortTypeName sorts entries by name
	SortTypeName SortType = "name"
	// SortTypeSize sorts entries by size
	SortTypeSize SortType = "size"
	// SortTypeDate sorts entries by date modified
	SortTypeDate SortType = "dateModified"
	// SortTypeExtension sorts entries by extension
	SortTypeExtension SortType = "extension"
	// SortTypeDirFirst sorts directories first
	SortTypeDirFirst SortType = "dirFirst"
)

// String returns the string representation of the sort type
func (s SortType) String() string {
	return string(s)
}
