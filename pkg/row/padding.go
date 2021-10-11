package row

import (
	"strings"
	"unicode/utf8"
)

// Right right-pads the string with pad up to len runes
func Right(str string, length int, pad string) string {
	count := length - utf8.RuneCountInString(str)
	if count <= 0 {
		count = 1
	}

	return str + strings.Repeat(pad, count)
}
