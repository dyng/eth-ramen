package util

import (
	"github.com/gdamore/tcell/v2"
)

const (
	NAValue = "n/a"

	EmptyValue = ""
)

type (
	KeyHandler func(*tcell.EventKey)

	KeyMap struct {
		Key         tcell.Key
		Shortcut    string
		Description string
		Handler     KeyHandler
	}

	KeyMaps []KeyMap
)

func (km KeyMaps) Add(another KeyMaps) KeyMaps {
	return append(km, another...)
}

func (km KeyMaps) FindHandler(key tcell.Key) (KeyHandler, bool) {
	for _, keymap := range km {
		if keymap.Key == key {
			return keymap.Handler, true
		}
	}
	return nil, false
}

// AsKey converts rune to keyboard key.
func AsKey(evt *tcell.EventKey) tcell.Key {
	if evt.Key() != tcell.KeyRune {
		return evt.Key()
	}
	return tcell.Key(evt.Rune())
}
