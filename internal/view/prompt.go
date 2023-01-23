package view

import (
	"github.com/dyng/ramen/internal/view/style"
	"github.com/ethereum/go-ethereum/log"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type PromptDialog struct {
	*tview.InputField
	display bool

	app *App
}

func NewPromptDialog(app *App) *PromptDialog {
	prompt := &PromptDialog{
		display:    false,
		app:        app,
	}

	// setup layout
	prompt.initLayout()

	return prompt
}

func (d *PromptDialog) initLayout() {
	s := d.app.config.Style()

	input := tview.NewInputField()
	input.SetFieldWidth(80)
	input.SetBorder(true)
	input.SetBorderColor(s.PromptBorderColor)
	input.SetTitle(style.Padding("Address"))
	input.SetTitleColor(s.FgColor)
	input.SetLabel("> ")
	input.SetLabelColor(s.InputFieldLableColor)
	input.SetFieldBackgroundColor(s.PromptBgColor)
	input.SetDoneFunc(d.handleKey)
	d.InputField = input
}

func (d *PromptDialog) handleKey(key tcell.Key) {
	switch key {
	case tcell.KeyEnter:
		address := d.GetText()
		if address != "" {
			account, err := d.app.service.GetAccount(address)
			if err != nil {
				// TODO: notify error
				log.Error("Failed to fetch account of given address", "address", address, "error", err)
			} else {
				d.app.root.HidePrompt()
				d.app.root.ShowAccountPage(account)
			}
		}
	case tcell.KeyEsc:
		d.app.root.HidePrompt()
	}
}

func (d *PromptDialog) Clear() {
	d.InputField.SetText("")
}

func (d *PromptDialog) Display(display bool) {
	d.display = display
}

func (d *PromptDialog) IsDisplay() bool {
	return d.display
}

// Draw implements tview.Primitive
func (d *PromptDialog) Draw(screen tcell.Screen) {
	if d.display {
		d.InputField.Draw(screen)
	}
}

// SetRect implements tview.SetRect
func (d *PromptDialog) SetRect(x int, y int, width int, height int) {
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

func (d *PromptDialog) inputSize() (int, int) {
	width := len(d.GetLabel()) + d.GetFieldWidth()
	height := d.GetFieldHeight() + 2
	return width, height
}
