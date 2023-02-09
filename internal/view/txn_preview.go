package view

import (
	"github.com/dyng/ramen/internal/view/util"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	// txnPreviewDialogMinHeight is the minimum height of the transaction preview dialog.
	txnPreviewDialogMinHeight = 20
	// txnPreviewDialogMinWidth is the minimum width of the transaction preview dialog.
	txnPreviewDialogMinWidth = 50
)

type TxnPreviewDialog struct {
	*TransactionDetail
	app       *App
	display   bool
	lastFocus tview.Primitive
}

func NewTxnPreviewDialog(app *App) *TxnPreviewDialog {
	d := &TxnPreviewDialog{
		TransactionDetail: NewTransactionDetail(app),
		app:               app,
		display:           false,
	}

	// setup keymap
	d.initKeymap()

	return d
}

func (d *TxnPreviewDialog) Show() {
	if !d.display {
		// save last focused element
		d.lastFocus = d.app.GetFocus()

		d.Display(true)
		d.app.SetFocus(d)
	}
}

func (d *TxnPreviewDialog) Hide() {
	if d.display {
		d.Display(false)
		d.app.SetFocus(d.lastFocus)
	}
}

func (d *TxnPreviewDialog) initKeymap() {
	InitKeymap(d, d.app)
}

// KeyMaps implements KeymapPrimitive
func (d *TxnPreviewDialog) KeyMaps() util.KeyMaps {
	keymaps := make(util.KeyMaps, 0)
	keymaps = append(keymaps, util.NewSimpleKey(tcell.KeyEsc, d.Hide))
	keymaps = append(keymaps, util.NewSimpleKey(tcell.KeyEnter, d.Hide))
	keymaps = append(keymaps, util.NewSimpleKey(util.KeySpace, d.Hide))
	keymaps = append(keymaps, util.NewSimpleKey(util.KeyF, func() {
		d.Hide()
		d.ViewSender()
	}))
	keymaps = append(keymaps, util.NewSimpleKey(util.KeyT, func() {
		d.Hide()
		d.ViewReceiver()
	}))
	return keymaps
}

func (d *TxnPreviewDialog) Display(display bool) {
	d.display = display
}

func (d *TxnPreviewDialog) IsDisplay() bool {
	return d.display
}

// Draw implements tview.Primitive
func (d *TxnPreviewDialog) Draw(screen tcell.Screen) {
	if d.display {
		d.TransactionDetail.Draw(screen)
	}
}

func (d *TxnPreviewDialog) SetCentral(x int, y int, width int, height int) {
	dialogWidth := width - width/3
	dialogHeight := height / 2
	if dialogWidth < txnPreviewDialogMinWidth {
		dialogWidth = txnPreviewDialogMinWidth
	}
	if dialogHeight < notificationMinHeight {
		dialogHeight = notificationMinHeight
	}
	dialogX := x + ((width - dialogWidth) / 2)
	dialogY := y + ((height - dialogHeight) / 2)
	d.TransactionDetail.SetRect(dialogX, dialogY, dialogWidth, dialogHeight)
}
