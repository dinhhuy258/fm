package row

import (
	"strings"
	"unicode/utf8"
)

// Right right-pads the string with pad up to len runes
func Right(str string, length int, pad string) string {
	return str + strings.Repeat(pad, length-utf8.RuneCountInString(str))
}
