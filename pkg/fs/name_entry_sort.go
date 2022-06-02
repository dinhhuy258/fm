package fs

import (
	"sort"
)

// nameEntrySort is an implementation of entrySort using sorting algorithm to sort entries by file size
type nameEntrySort struct{}

// sort by file name
func (nameEntrySort) sort(entries []IEntry, reverse bool, ignoreCase bool, ignoreDiacritics bool) {
	sort.Slice(entries, func(i, j int) bool {
		name1 := normalize(entries[i].GetName(), ignoreCase, ignoreDiacritics)
		name2 := normalize(entries[j].GetName(), ignoreCase, ignoreDiacritics)

		s := name1 < name2
		if reverse {
			s = !s
		}

		return s
	})
}
