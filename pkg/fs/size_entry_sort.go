package fs

import (
	"sort"
)

// sizeEntrySor is an implementation of entrySort using sorting algorithm to sort entries by file size
type sizeEntrySort struct{}

// sort by file size
func (sizeEntrySort) sort(entries []IEntry, reverse bool, ignoreCase bool, ignoreDiacritics bool) {
	sort.Slice(entries, func(i, j int) bool {
		s := entries[i].GetSize() > entries[j].GetSize()
		if reverse {
			s = !s
		}

		return s
	})
}
