package fs

import (
	"sort"
)

// entryTypeSort represents struct for sorting algorithm by dir
type dirFirstEntrySort struct{}

// sort the entries according to dir first algorithm
func (_ dirFirstEntrySort) sort(entries []IEntry) {
	sort.Slice(entries, func(i, j int) bool {
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
	})
}
