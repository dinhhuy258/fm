package fs

import (
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// entrySort an interface for entry sorting algorithm
type entrySort interface {
	sort([]IEntry, bool, bool, bool)
}

// sortType is the type of sorting algorithm
type sortType string

const (
	DirFirst     sortType = "dirFirst"
	DateModified sortType = "dateModified"
	Name         sortType = "name"
	Size         sortType = "size"
	Extension    sortType = "extension"
)

// entrySortFactories contains list of supported entry sort algorithm
var entrySortFactories = map[sortType]entrySort{
	DirFirst:     dirFirstEntrySort{},
	DateModified: dateModifiedEntrySort{},
	Name:         nameEntrySort{},
	Size:         sizeEntrySort{},
	Extension:    extensionEntrySort{},
}

// getEntrySort returns entry sort algorithm for the given sortType
func getEntrySort(t sortType) entrySort {
	entrySort, hashEntrySort := entrySortFactories[t]
	if !hashEntrySort {
		// fallback to DirFirst
		entrySort = entrySortFactories[DirFirst]
	}

	return entrySort
}

// normalize the given string
func normalize(s string, ignoreCase, ignoreDiacritics bool) string {
	if ignoreCase {
		s = strings.ToLower(s)
	}

	if ignoreDiacritics {
		s = removeDiacritics(s)
	}

	return s
}

// removeDiacritics from the given string
func removeDiacritics(str string) string {
	isMn := func(r rune) bool {
		return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
	}

	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	s, _, err := transform.String(t, str)
	if err != nil {
		return str
	}

	return s
}
