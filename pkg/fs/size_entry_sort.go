package fs

import (
	"sort"
)

// sizeEntrySor is an implementation of entrySort using sorting algorithm to sort entries by file size
type sizeEntrySort struct{}

// sort by file size
func (sizeEntrySort) sort(entries []IEntry, reverse bool, ignoreCase bool, ignoreDiacritics bool) {
	sort.Slice(entries, func(i, j int) bool {
		entry1 := entries[i]
		entry2 := entries[j]

		if !entry1.IsDirectory() && !entry2.IsDirectory() {
			// if both entries are files
			s := entry1.GetSize() > entry2.GetSize()
			if reverse {
				s = !s
			}

			return s
		}

		if entry1.IsDirectory() && entry2.IsDirectory() {
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
	})
}
