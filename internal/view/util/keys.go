package util

import (
	"github.com/gdamore/tcell/v2"
)

const (
	NAValue = "[dimgray]n/a[-]"

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

// NewSimpleKey creates a simple keymap with only key and handler.
func NewSimpleKey(key tcell.Key, handler func()) KeyMap {
	return KeyMap{
		Key:     key,
		Handler: func(*tcell.EventKey) { handler() },
	}
}

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
