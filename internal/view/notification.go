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

	title string
	text  string
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

func (n *Notification) initLayout() {
	s := n.app.config.Style()

	tv := tview.NewTextView()
	tv.SetBorder(true)
	tv.SetBorderColor(s.SecondaryBorderColor)
	tv.SetWrap(true)
	n.TextView = tv
}

func (n *Notification) initKeymap() {
	n.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := util.AsKey(event)
		switch key {
		case tcell.KeyEsc, tcell.KeyEnter, util.KeySpace:
			n.app.root.HideNotification()
			return nil
		default:
			return event
		}
	})
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

// SetRect implements tview.SetRect
func (n *Notification) SetRect(x int, y int, width int, height int) {
	dialogWidth := width - width/3
	dialogHeight := height / 4
	if dialogHeight < 15 {
		dialogHeight = 15
	}
	dialogX := x + ((width - dialogWidth) / 2)
	dialogY := y + ((height - dialogHeight) / 2)
	n.TextView.SetRect(dialogX, dialogY, dialogWidth, dialogHeight)
}
