package view

import (
	"github.com/dyng/ramen/internal/view/style"
	"github.com/dyng/ramen/internal/view/util"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Notification struct {
	*tview.TextView
	app     *App
	display bool

	title     string
	text      string
	lastFocus tview.Primitive
}

func NewNotification(app *App) *Notification {
	n := &Notification{
		app:     app,
		display: false,
	}

	// setup layout
	n.initLayout()

	// setup keymap
	n.initKeymap()

	return n
}

func (n *Notification) SetContent(title string, text string) {
	n.title = title
	n.text = text
	n.refresh()
}

func (n *Notification) Show() {
	// save last focused element
	n.lastFocus = n.app.GetFocus()

	n.Display(true)
	n.app.SetFocus(n)
}

func (n *Notification) Hide() {
	n.Display(false)
	n.app.SetFocus(n.lastFocus)
}

func (n *Notification) initLayout() {
	s := n.app.config.Style()

	tv := tview.NewTextView()
	tv.SetBorder(true)
	tv.SetBorderColor(s.BorderColor2)
	tv.SetWrap(true)
	n.TextView = tv
}

func (n *Notification) initKeymap() {
	InitKeymap(n, n.app)
}

// KeyMaps implements KeymapPrimitive
func (n *Notification) KeyMaps() util.KeyMaps {
	keymaps := make(util.KeyMaps, 0)
	keymaps = append(keymaps, util.NewSimpleKey(tcell.KeyEsc, n.Hide))
	keymaps = append(keymaps, util.NewSimpleKey(tcell.KeyEnter, n.Hide))
	keymaps = append(keymaps, util.NewSimpleKey(util.KeySpace, n.Hide))
	return keymaps
}

func (n *Notification) refresh() {
	n.SetTitle(style.BoldPadding(n.title))
	n.SetText(n.text)
}

func (n *Notification) Display(display bool) {
	n.display = display
}

func (n *Notification) IsDisplay() bool {
	return n.display
}

// Draw implements tview.Primitive
func (n *Notification) Draw(screen tcell.Screen) {
	if n.display {
		n.TextView.Draw(screen)
	}
}

func (n *Notification) SetCentral(x int, y int, width int, height int) {
	dialogWidth := width - width/3
	dialogHeight := height / 4
	if dialogHeight < 15 {
		dialogHeight = 15
	}
	dialogX := x + ((width - dialogWidth) / 2)
	dialogY := y + ((height - dialogHeight) / 2)
	n.TextView.SetRect(dialogX, dialogY, dialogWidth, dialogHeight)
}
