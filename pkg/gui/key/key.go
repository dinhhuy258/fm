package key

import (
	"fmt"
	"log"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/dinhhuy258/gocui"
)

// Key represent gocui key
type Key struct {
	Key gocui.Key
	Ch  rune
	Mod gocui.Modifier
}

// GetKeyDisplay returns the display name of the key
func GetKeyDisplay(key Key) string {
	keyInt := 0

	if key.Ch != 0 && key.Mod == gocui.ModNone && unicode.IsPrint(key.Ch) {
		keyInt = int(key.Ch)
	} else {
		value, ok := keyMapReversed[key.Key]
		if ok {
			return value
		}

		keyInt = int(key.Ch)
	}

	return fmt.Sprintf("%c", keyInt)
}

// GetKey returns the key from the display name
func GetKey(key string) Key {
	runeCount := utf8.RuneCountInString(key)
	if runeCount > 1 {
		binding, hasKey := keymap[strings.ToLower(key)]
		if !hasKey {
			log.Fatalf("Unrecognized key %s for keybinding", strings.ToLower(key))
		} else {
			return binding
		}
	} else if runeCount == 1 {
		return Key{
			Key: 0,
			Ch:  []rune(key)[0],
			Mod: gocui.ModNone,
		}
	}

	log.Fatal("Key empty for keybinding: " + strings.ToLower(key))

	return Key{
		Key: 0,
		Ch:  0,
		Mod: gocui.ModNone,
	}
}
