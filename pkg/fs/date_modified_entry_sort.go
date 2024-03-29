package fs

import (
	"sort"
)

// dateModifiedEntrySort is an implementation of entrySort using sorting algorithm to sort entries
// by last modified time
type dateModifiedEntrySort struct{}

// sort by last modified time
func (dateModifiedEntrySort) sort(entries []IEntry, reverse bool, _ bool, _ bool) {
	sort.Slice(entries, func(i, j int) bool {
		s := entries[i].GetChangeTime().After(entries[j].GetChangeTime())
		if reverse {
			s = !s
		}

		return s
	})
}
