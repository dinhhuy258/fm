package fs

import (
	"sort"
)

// nameEntrySort is an implementation of entrySort using sorting algorithm to sort entries by file size
type nameEntrySort struct{}

// sort by file name
func (nameEntrySort) sort(entries []IEntry, reverse bool) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].GetName() < entries[j].GetName()
	})
}
