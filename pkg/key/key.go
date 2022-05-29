package key

import (
	"fmt"
	"log"
	"strings"
	"unicode/utf8"

	"github.com/dinhhuy258/gocui"
)

type Key interface{} // FIXME: find out how to get `gocui.Key | rune`

// GetKeyDisplay returns the display name of the key
func GetKeyDisplay(key Key) string {
	keyInt := 0

	switch key := key.(type) {
	case rune:
		keyInt = int(key)
	case gocui.Key:
		value, ok := keyMapReversed[key]
		if ok {
			return value
		}

		keyInt = int(key)
	}

	return fmt.Sprintf("%c", keyInt)
}

// GetKey returns the key from the display name
func GetKey(key string) Key {
	runeCount := utf8.RuneCountInString(key)
	if runeCount > 1 {
		binding := keymap[strings.ToLower(key)]
		if binding == nil {
			log.Fatalf("Unrecognized key %s for keybinding", strings.ToLower(key))
		} else {
			return binding
		}
	} else if runeCount == 1 {
		return []rune(key)[0]
	}

	log.Fatal("Key empty for keybinding: " + strings.ToLower(key))

	return nil
}
