package view

import (
	"github.com/dyng/ramen/internal/view/format"
	"github.com/dyng/ramen/internal/view/style"
	"github.com/ethereum/go-ethereum/log"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type SignInDialog struct {
	*tview.InputField
	app       *App
	display   bool
	lastFocus tview.Primitive
}

func NewSignInDialog(app *App) *SignInDialog {
	d := &SignInDialog{
		app:     app,
		display: false,
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
		// close dialog at first
		d.Hide()

		privateKey := d.GetText()
		if privateKey != "" {
			signer, err := d.app.service.GetSigner(privateKey)
			if err != nil {
				log.Error("Failed to create signer", "error", err)
				d.app.root.NotifyError(format.FineErrorMessage("Failed to create signer", err))
			} else {
				d.app.root.SignIn(signer)
			}
		}
	case tcell.KeyEsc:
		d.Hide()
	}
}

func (d *SignInDialog) Show() {
	// save last focused element
	d.lastFocus = d.app.GetFocus()

	d.Display(true)
	d.app.SetFocus(d)
}

func (d *SignInDialog) Hide() {
	d.Display(false)
	d.app.SetFocus(d.lastFocus)
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
