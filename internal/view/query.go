package view

import (
	"github.com/dyng/ramen/internal/view/format"
	"github.com/dyng/ramen/internal/view/style"
	"github.com/dyng/ramen/internal/view/util"
	"github.com/ethereum/go-ethereum/log"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type QueryDialog struct {
	*tview.InputField
	app     *App
	display bool
	spinner *util.Spinner
}

func NewQueryDialog(app *App) *QueryDialog {
	query := &QueryDialog{
		app:        app,
		display:    false,
		spinner:    util.NewSpinner(app.Application),
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
		// start account query
		d.setSpinnerRect()
		d.spinner.StartAndShow()

		address := d.GetText()
		query := func() {
			if address != "" {
				account, err := d.app.service.GetAccount(address)
				d.app.QueueUpdateDraw(func() {
					d.spinner.StopAndHide()
					d.app.root.HideQueryDialog()

					if err != nil {
						log.Error("Failed to fetch account of given address",
							"address", address, "error", err)
						d.app.root.NotifyError(format.FineErrorMessage(
							"Failed to fetch account of address %s", address, err))
					} else {
						d.app.root.ShowAccountPage(account)
					}
				})
			}
		}
		go query()
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
	d.spinner.Draw(screen)
}

func (d *QueryDialog) SetCentral(x int, y int, width int, height int) {
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

func (d *QueryDialog) setSpinnerRect() {
	x, y, _, _ := d.GetInnerRect()
	sx := x + len(d.GetText()) + 2
	d.spinner.SetRect(sx, y, 0, 0)
}
