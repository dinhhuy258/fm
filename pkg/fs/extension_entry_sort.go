package fs

import (
	"sort"
)

// extensionEntrySort is an implementation of entrySort using sorting algorithm to sort entries by extension
type extensionEntrySort struct{}

// sort by extension
func (extensionEntrySort) sort(entries []IEntry, reverse bool, ignoreCase bool, ignoreDiacritics bool) {
	sort.Slice(entries, func(i, j int) bool {
		s := func() bool {
			ext1 := normalize(entries[i].GetExt(), ignoreCase, ignoreDiacritics)
			ext2 := normalize(entries[j].GetExt(), ignoreCase, ignoreDiacritics)
			name1 := normalize(entries[i].GetName(), ignoreCase, ignoreDiacritics)
			name2 := normalize(entries[j].GetName(), ignoreCase, ignoreDiacritics)

			// if the extension could not be determined (directories, files without)
			// use a zero byte so that these files can be ranked higher
			if ext1 == "" {
				ext1 = "\x00"
			}
			if ext2 == "" {
				ext2 = "\x00"
			}

			return ext1 < ext2 || ext1 == ext2 && name1 < name2
		}()

		if reverse {
			s = !s
		}

		return s
	})
}
