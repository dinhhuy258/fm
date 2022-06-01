package fs

import (
	"sort"
)

// entryTypeSort is an implementation of entrySort using sorting algorithm to sort entries by dir
// first
type dirFirstEntrySort struct{}

// sort by dir first
func (dirFirstEntrySort) sort(entries []IEntry, reverse bool) {
	sort.Slice(entries, func(i, j int) bool {
		s := func() bool {
			firstEntry := entries[i]
			secondEntry := entries[j]

			if (firstEntry.IsDirectory() && secondEntry.IsDirectory()) ||
				(!firstEntry.IsDirectory() && !secondEntry.IsDirectory()) {
				// If 2 entries are both directory or file, then sort by their name
				return firstEntry.GetName() < secondEntry.GetName()
			}

			if firstEntry.IsDirectory() {
				return true
			}

			// Second entry is directory
			return false
		}()

		if reverse {
			s = !s
		}

		return s
	})
}
