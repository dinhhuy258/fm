package key

import (
	"fmt"
	"log"
	"strings"
	"unicode"

	"github.com/dinhhuy258/gocui"
	"github.com/rivo/uniseg"
)

// Key represent gocui key
type Key struct {
	Ch  rune
	Key gocui.Key
	Mod gocui.Modifier
}

// IsDigit check if the key is digit
func (k Key) IsDigit() bool {
	return k.Mod == gocui.ModNone && unicode.IsDigit(k.Ch)
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
	runeCount := uniseg.GraphemeClusterCount(key)
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

	// Not reachable
	return Key{
		Key: 0,
		Ch:  0,
		Mod: gocui.ModNone,
	}
}
