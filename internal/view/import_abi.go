package view

import (
	"github.com/dyng/ramen/internal/view/format"
	"github.com/dyng/ramen/internal/view/style"
	"github.com/dyng/ramen/internal/view/util"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ImportABIDialog struct {
	*tview.Flex
	app     *App
	display bool

	input *tview.TextArea
	button *tview.Button
}

func NewImportABIDialog(app *App) *ImportABIDialog {
	d := &ImportABIDialog{
		app:     app,
		display: false,
	}

	// setup layout
	d.initLayout()

	// setup keymap
	d.initKeymap()

	return d
}

func (d *ImportABIDialog) initLayout() {
	s := d.app.config.Style()

	// description
	desc := tview.NewTextView()
	desc.SetWrap(true)
	desc.SetTextAlign(tview.AlignCenter)
	desc.SetBorderPadding(0, 0, 1, 1)
	desc.SetText("Cannot find ABI for this contract. But you can upload an ABI json instead.\nGenerate ABI json by solc command: `solc filename.sol --abi`.")

	// textarea
	input := tview.NewTextArea()
	input.SetBorder(true)
	input.SetBorderColor(s.SecondaryBorderColor)
	input.SetWrap(true)
	d.input = input

	// buttons
	buttons := tview.NewForm()
	buttons.SetButtonsAlign(tview.AlignRight)
	buttons.SetButtonBackgroundColor(s.ButtonBgColor)
	buttons.AddButton("Import", d.doImport)
	d.button = buttons.GetButton(0)

	// flex
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetBorder(true)
	flex.SetBorderColor(s.DialogBorderColor)
	flex.SetTitle(style.BoldPadding("Import ABI"))
	flex.AddItem(desc, 0, 2, false)
	flex.AddItem(input, 0, 8, true)
	flex.AddItem(buttons, 3, 0, false)

	d.Flex = flex
}

func (d *ImportABIDialog) initKeymap() {
	d.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := util.AsKey(event)
		switch key {
		case tcell.KeyEsc:
			d.app.root.account.HideImportABIDialog()
			return nil
		case tcell.KeyTab:
			d.focusNext()
			return nil
		default:
			return event
		}
	})
}

func (d *ImportABIDialog) focusNext() {
	if d.input.HasFocus() {
		d.app.SetFocus(d.button)
		return
	}

	if d.button.HasFocus() {
		d.app.SetFocus(d.input)
		return
	}
}

func (d *ImportABIDialog) doImport() {
	account := d.app.root.account

	// hide dialog
	account.HideImportABIDialog()
	
	// read and parse abi json
	err := account.contract.ImportABI(d.input.GetText())
	if err != nil {
		d.app.root.NotifyError(format.FineErrorMessage("Cannot import ABI json", err))
		return
	}

	// show callMethod dialog
	account.methodCall.refresh()
	account.ShowMethodCallDialog()
}

func (d *ImportABIDialog) Clear() {
	d.input.SetText("", true)
}

func (d *ImportABIDialog) Display(display bool) {
	d.display = display
}

func (d *ImportABIDialog) IsDisplay() bool {
	return d.display
}

// Draw implements tview.Primitive
func (d *ImportABIDialog) Draw(screen tcell.Screen) {
	if d.display {
		d.Flex.Draw(screen)
	}
}

// SetRect implements tview.SetRect
func (d *ImportABIDialog) SetRect(x int, y int, width int, height int) {
	dialogWidth := width - width/2
	dialogHeight := height - height/2
	if dialogHeight < 10 {
		dialogHeight = 10
	}
	dialogX := x + ((width - dialogWidth) / 2)
	dialogY := y + ((height - dialogHeight) / 2)
	d.Flex.SetRect(dialogX, dialogY, dialogWidth, dialogHeight)
}
