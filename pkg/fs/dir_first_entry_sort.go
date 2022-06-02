package fs

import (
	"sort"
)

// entryTypeSort is an implementation of entrySort using sorting algorithm to sort entries by dir
// first
type dirFirstEntrySort struct{}

// sort by dir first
func (dirFirstEntrySort) sort(entries []IEntry, reverse bool, ignoreCase bool, ignoreDiacritics bool) {
	sort.Slice(entries, func(i, j int) bool {
		s := func() bool {
			entry1 := entries[i]
			entry2 := entries[j]

			if (entry1.IsDirectory() && entry2.IsDirectory()) ||
				(!entry1.IsDirectory() && !entry2.IsDirectory()) {
				// If 2 entries are both directory or file, then sort by their name
				name1 := normalize(entry1.GetName(), ignoreCase, ignoreDiacritics)
				name2 := normalize(entry2.GetName(), ignoreCase, ignoreDiacritics)

				return name1 < name2
			}

			if entry1.IsDirectory() {
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
