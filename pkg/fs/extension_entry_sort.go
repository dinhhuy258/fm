package fs

import (
	"sort"
)

// extensionEntrySort is an implementation of entrySort using sorting algorithm to sort entries by extension
type extensionEntrySort struct{}

// sort by extension
func (extensionEntrySort) sort(entries []IEntry, reverse bool) {
	sort.Slice(entries, func(i, j int) bool {
		ext1 := entries[i].GetExt()
		ext2 := entries[j].GetExt()
		name1 := entries[i].GetName()
		name2 := entries[j].GetName()

		// if the extension could not be determined (directories, files without)
		// use a zero byte so that these files can be ranked higher
		if ext1 == "" {
			ext1 = "\x00"
		}
		if ext2 == "" {
			ext2 = "\x00"
		}

		return ext1 < ext2 || ext1 == ext2 && name1 < name2
	})
}
