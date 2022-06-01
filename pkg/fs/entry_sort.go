package fs

// entrySort an interface for entry sorting algorithm
type entrySort interface {
	sort([]IEntry, bool)
}

// sortType is the type of sorting algorithm
type sortType string

const (
	DirFirst     = "dirFirst"
	DateModified = "dateModified"
	Name         = "name"
	Size         = "size"
	Extension    = "extension"
)

// entrySortFactories contains list of supported entry sort algorithm
var entrySortFactories = map[sortType]entrySort{
	DirFirst:     dirFirstEntrySort{},
	DateModified: dateModifiedEntrySort{},
	Name:         nameEntrySort{},
	Size:         sizeEntrySort{},
	Extension:    extensionEntrySort{},
}

// getEntrySort returns entry sort algorithm for the given sortType
func getEntrySort(t sortType) entrySort {
	entrySort, hashEntrySort := entrySortFactories[t]
	if !hashEntrySort {
		// fallback to DirFirst
		entrySort = entrySortFactories[DirFirst]
	}

	return entrySort
}
