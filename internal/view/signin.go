package view

import (
	"github.com/dyng/ramen/internal/view/format"
	"github.com/dyng/ramen/internal/view/style"
	"github.com/dyng/ramen/internal/view/util"
	"github.com/ethereum/go-ethereum/log"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type SignInDialog struct {
	*tview.InputField
	app       *App
	display   bool
	lastFocus tview.Primitive
	spinner   *util.Spinner
}

func NewSignInDialog(app *App) *SignInDialog {
	d := &SignInDialog{
		app:     app,
		display: false,
		spinner: util.NewSpinner(app.Application),
	}

	// setup layout
	d.initLayout()

	return d
}

func (d *SignInDialog) initLayout() {
	s := d.app.config.Style()

	input := tview.NewInputField()
	input.SetFieldWidth(80)
	input.SetBorder(true)
	input.SetBorderColor(s.DialogBorderColor)
	input.SetTitle(style.Padding("Private Key"))
	input.SetTitleColor(s.FgColor)
	input.SetLabel(" ")
	input.SetMaskCharacter('*')
	input.SetFieldBackgroundColor(s.DialogBgColor)
	input.SetDoneFunc(d.handleKey)
	d.InputField = input
}

func (d *SignInDialog) handleKey(key tcell.Key) {
	switch key {
	case tcell.KeyEnter:
		// start spinner
		d.Loading()

		privateKey := d.GetText()
		if privateKey == "" {
			return
		}

		go func() {
			signer, err := d.app.service.GetSigner(privateKey)
			signer.UpdateBalance() // populate balance cache
			d.app.QueueUpdateDraw(func() {
				if err != nil {
					d.Finished()
					log.Error("Failed to create signer", "error", err)
					d.app.root.NotifyError(format.FineErrorMessage("Failed to create signer", err))
				} else {
					d.app.root.SignIn(signer)
					d.Finished()
				}
			})
		}()
	case tcell.KeyEsc:
		d.Hide()
	}
}

func (d *SignInDialog) Show() {
	if !d.display {
		// save last focused element
		d.lastFocus = d.app.GetFocus()

		d.Display(true)
		d.app.SetFocus(d)
	}
}

func (d *SignInDialog) Hide() {
	if d.display {
		d.Display(false)
		d.app.SetFocus(d.lastFocus)
	}
}

// Loading will set the location of spinner and show it
func (d *SignInDialog) Loading() {
	d.setSpinnerRect()
	d.spinner.StartAndShow()
}

// Finished will stop and hide spinner, as well as close current dialog
func (d *SignInDialog) Finished() {
	d.spinner.StopAndHide()
	d.Hide()
}

func (d *SignInDialog) Clear() {
	d.InputField.SetText("")
}

func (d *SignInDialog) Display(display bool) {
	d.display = display
}

func (d *SignInDialog) IsDisplay() bool {
	return d.display
}

// Draw implements tview.Primitive
func (d *SignInDialog) Draw(screen tcell.Screen) {
	if d.display {
		d.InputField.Draw(screen)
	}
	d.spinner.Draw(screen)
}

func (d *SignInDialog) SetCentral(x int, y int, width int, height int) {
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

func (d *SignInDialog) inputSize() (int, int) {
	width := len(d.GetLabel()) + d.GetFieldWidth()
	height := d.GetFieldHeight() + 2
	return width, height
}

func (d *SignInDialog) setSpinnerRect() {
	x, y, _, _ := d.GetInnerRect()
	sx := x + len(d.GetText()) + 1
	d.spinner.SetRect(sx, y, 0, 0)
}
