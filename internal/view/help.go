package view

import (
	"fmt"

	"github.com/dyng/ramen/internal/view/style"
	"github.com/dyng/ramen/internal/view/util"
	"github.com/rivo/tview"
)

type Help struct {
	*tview.Table
	app *App

	keymaps util.KeyMaps
}

func NewHelp(app *App) *Help {
	help := &Help{
		Table: tview.NewTable(),
		app:   app,
	}

	return help
}

func (h *Help) SetKeyMaps(keymaps util.KeyMaps) {
	h.keymaps = keymaps
	h.refresh()
}

func (h *Help) AddKeyMaps(keymaps util.KeyMaps) {
	h.keymaps.Add(keymaps)
	h.refresh()
}

func (h *Help) Clear() {
	for i := h.GetRowCount() - 1; i > 0; i-- {
		h.RemoveRow(i)
	}
}

func (h *Help) refresh() {
	// clear previous content at first
	h.Clear()

	s := h.app.config.Style()
	row, col := 0, 0
	for _, keymap := range h.keymaps {
		if row >= style.HeaderHeight {
			row = 0
			col += 2
		}

		short := fmt.Sprintf("<%s>", keymap.Shortcut)
		desc := keymap.Description
		sec := util.NewSectionWithColor(short, s.HelpKeyColor, desc, s.FgColor)
		sec.AddToTable(h.Table, row, col)

		row += 1
	}
}
