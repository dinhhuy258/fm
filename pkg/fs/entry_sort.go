package fs

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"github.com/dinhhuy258/fm/pkg/types"
)

// entrySort an interface for entry sorting algorithm
type entrySort interface {
	sort([]IEntry, bool, bool, bool)
}

// entrySortFactories contains list of supported entry sort algorithm
var entrySortFactories = map[types.SortType]entrySort{
	types.SortTypeDirFirst:  dirFirstEntrySort{},
	types.SortTypeDate:      dateModifiedEntrySort{},
	types.SortTypeName:      nameEntrySort{},
	types.SortTypeSize:      sizeEntrySort{},
	types.SortTypeExtension: extensionEntrySort{},
}

// getEntrySort returns entry sort algorithm for the given sortType
func getEntrySort(t types.SortType) entrySort {
	entrySort, hashEntrySort := entrySortFactories[t]
	if !hashEntrySort {
		// fallback to DirFirst
		entrySort = entrySortFactories[types.SortTypeDirFirst]
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
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

	s, _, err := transform.String(t, str)
	if err != nil {
		return str
	}

	return s
}
