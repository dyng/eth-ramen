package view

import (
	"github.com/dyng/ramen/internal/view/format"
	"github.com/dyng/ramen/internal/view/style"
	"github.com/ethereum/go-ethereum/log"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type QueryDialog struct {
	*tview.InputField
	app *App
	display bool
}

func NewQueryDialog(app *App) *QueryDialog {
	query := &QueryDialog{
		display:    false,
		app:        app,
	}

	// setup layout
	query.initLayout()

	return query
}

func (d *QueryDialog) initLayout() {
	s := d.app.config.Style()

	input := tview.NewInputField()
	input.SetFieldWidth(80)
	input.SetBorder(true)
	input.SetBorderColor(s.DialogBorderColor)
	input.SetTitle(style.Padding("Address"))
	input.SetTitleColor(s.FgColor)
	input.SetLabel("> ")
	input.SetLabelColor(s.InputFieldLableColor)
	input.SetFieldBackgroundColor(s.DialogBgColor)
	input.SetDoneFunc(d.handleKey)
	d.InputField = input
}

func (d *QueryDialog) handleKey(key tcell.Key) {
	switch key {
	case tcell.KeyEnter:
		address := d.GetText()
		if address != "" {
			account, err := d.app.service.GetAccount(address)
			if err != nil {
				log.Error("Failed to fetch account of given address", "address", address, "error", err)
				d.app.root.NotifyError(format.FineErrorMessage(
					"Failed to fetch account of address %s", address, err))
			} else {
				d.app.root.HideQueryDialog()
				d.app.root.ShowAccountPage(account)
			}
		}
	case tcell.KeyEsc:
		d.app.root.HideQueryDialog()
	}
}

func (d *QueryDialog) Clear() {
	d.InputField.SetText("")
}

func (d *QueryDialog) Display(display bool) {
	d.display = display
}

func (d *QueryDialog) IsDisplay() bool {
	return d.display
}

// Draw implements tview.Primitive
func (d *QueryDialog) Draw(screen tcell.Screen) {
	if d.display {
		d.InputField.Draw(screen)
	}
}

// SetRect implements tview.SetRect
func (d *QueryDialog) SetRect(x int, y int, width int, height int) {
	inputWidth, inputHeight := d.inputSize()
	if inputWidth > width-2 {
		inputWidth = width - 2
	}
	if inputHeight > height-2 {
		inputHeight = height
	}
	ws := (width - inputWidth) / 2
	hs := (height - inputHeight) / 2
	d.InputField.SetRect(x+ws, y+hs, inputWidth, inputHeight)
}

func (d *QueryDialog) inputSize() (int, int) {
	width := len(d.GetLabel()) + d.GetFieldWidth()
	height := d.GetFieldHeight() + 2
	return width, height
}
