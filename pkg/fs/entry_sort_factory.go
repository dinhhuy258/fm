package fs

type sortType int8

const (
	DirFirst sortType = iota
)

// entrySortFactories contains list of supported entry sort algorithm
var entrySortFactories = map[sortType]entrySort{
	DirFirst: dirFirstEntrySort{},
}

// getEntrySort respond the entrySort corresponding with the given sort type
func getEntrySort(t sortType) entrySort {
	entrySort, hashEntrySort := entrySortFactories[t]
	if !hashEntrySort {
		return entrySortFactories[DirFirst]
	}

	return entrySort
}
