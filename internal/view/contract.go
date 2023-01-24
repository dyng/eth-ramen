package view

import (
	"fmt"

	"github.com/dyng/ramen/internal/common/conv"
	"github.com/dyng/ramen/internal/service"
	"github.com/dyng/ramen/internal/view/format"
	"github.com/dyng/ramen/internal/view/style"
	"github.com/dyng/ramen/internal/view/util"
	"github.com/ethereum/go-ethereum/log"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type MethodCallDialog struct {
	*tview.Flex
	app     *App
	display bool

	methods  *tview.Table
	args     *tview.Form
	result   *tview.TextView
	focusIdx int
	spinner  *util.Spinner

	contract *service.Contract
}

func NewMethodCallDialog(app *App) *MethodCallDialog {
	mcd := &MethodCallDialog{
		app: app,
		spinner: util.NewSpinner(app.Application),
	}

	// setup layout
	mcd.initLayout()

	return mcd
}

func (d *MethodCallDialog) initLayout() {
	s := d.app.config.Style()

	// method list
	methods := tview.NewTable()
	methods.SetBorder(true)
	methods.SetBorderColor(s.MethNameBorderColor)
	methods.SetTitle(style.Padding("Method"))
	methods.SetSelectable(true, false)
	methods.SetSelectionChangedFunc(func(row, column int) {
		d.showArguments()
	})
	methods.SetSelectedFunc(func(row, column int) {
		if d.methodHasNoArgs() {
			d.callMethod()
		} else {
			d.focusNext()
		}
	})
	methods.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := util.AsKey(event)
		switch key {
		case tcell.KeyTAB:
			d.focusNext()
			return nil
		default:
			return event
		}
	})
	d.methods = methods

	// arguments form
	args := tview.NewForm()
	args.SetBorder(true)
	args.SetBorderColor(s.MethArgsBorderColor)
	args.SetTitle(style.Padding("Arguments"))
	args.SetLabelColor(s.InputFieldLableColor)
	args.SetFieldBackgroundColor(s.InputFieldBgColor)
	args.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := util.AsKey(event)
		switch key {
		case tcell.KeyEnter:
			if d.atLastFormItem() {
				d.callMethod()
			} else {
				d.focusNext()
			}
			return nil
		case tcell.KeyTAB:
			d.focusNext()
			return nil
		default:
			return event
		}
	})
	d.args = args

	top := tview.NewFlex().SetDirection(tview.FlexColumn)
	top.AddItem(methods, 0, 3, false)
	top.AddItem(args, 0, 7, true)

	// result panel
	result := tview.NewTextView()
	result.SetBorder(true)
	result.SetBorderColor(s.MethResultBorderColor)
	result.SetTitle(style.Padding("Result"))
	d.result = result

	whole := tview.NewFlex().SetDirection(tview.FlexRow)
	whole.AddItem(top, 0, 8, true)
	whole.AddItem(result, 0, 2, false)
	whole.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := util.AsKey(event)
		if key == tcell.KeyEscape {
			d.app.root.account.HideMethodCallDialog()
			return nil
		} else {
			return event
		}
	})

	d.Flex = whole
}

func (d *MethodCallDialog) SetContract(contract *service.Contract) {
	d.contract = contract
	d.refresh()
}

func (d *MethodCallDialog) Clear() {
	d.result.Clear()
}

func (d *MethodCallDialog) refresh() {
	d.args.Clear(true)
	d.result.Clear()
	d.showMethodList()
}

func (d *MethodCallDialog) showMethodList() {
	d.methods.Clear()

	s := d.app.config.Style()
	row := 0
	for name, method := range d.contract.GetABI().Methods {
		if method.IsConstant() {
			d.methods.SetCell(row, 0, tview.NewTableCell(name).SetTextColor(s.FgColor))
			row++
		}
	}
}

func (d *MethodCallDialog) showArguments() {
	d.args.Clear(true)

	methodName := d.methods.GetCell(d.methods.GetSelection()).Text
	method := d.contract.GetABI().Methods[methodName]
	for _, arg := range method.Inputs {
		argName := arg.Name
		if argName == "" {
			argName = "<unknown>"
		}
		d.args.AddInputField(argName, "", 999, nil, nil)
	}
}

func (d *MethodCallDialog) focusNext() {
	next := d.focusIdx + 1
	if next > d.args.GetFormItemCount() {
		next = 0
	}

	if next == 0 {
		d.app.SetFocus(d.methods)
	} else {
		formItem := d.args.GetFormItem(next - 1)
		d.app.SetFocus(formItem)
	}

	d.focusIdx = next
}

func (d *MethodCallDialog) atLastFormItem() bool {
	return d.focusIdx >= d.args.GetFormItemCount()
}

func (d *MethodCallDialog) methodHasNoArgs() bool {
	methodName := d.methods.GetCell(d.methods.GetSelection()).Text
	method := d.contract.GetABI().Methods[methodName]
	return len(method.Inputs) == 0
}

func (d *MethodCallDialog) callMethod() {
	methodName := d.methods.GetCell(d.methods.GetSelection()).Text
	method := d.contract.GetABI().Methods[methodName]

	args := make([]any, 0)
	for i := 0; i < d.args.GetFormItemCount(); i++ {
		item := d.args.GetFormItem(i).(*tview.InputField)
		arg := method.Inputs[i]
		val, err := conv.Unpack(arg.Type, item.GetText())
		if err == nil {
			args = append(args, val)
		} else {
			log.Error("Cannot unpack argument", "argument", arg, "input", item.GetText(), "error", err)
			d.app.root.NotifyError(format.FineErrorMessage(
				"Input type for argument '%s' is incorrect, should be '%s'.", arg.Name, arg.Type.String(), err))
			return
		}
	}

	// start calling...
	d.showSpinner()

	call := func() {
		res, err := d.contract.Call(methodName, args...)
		d.app.QueueUpdateDraw(func() {
			if err != nil {
				log.Error("Method call is failed", "name", methodName, "args", args, "error", err)
				d.app.root.NotifyError(format.FineErrorMessage("Cannot call contract method '%s'.", methodName, err))
			} else {
				d.result.SetText(fmt.Sprint(res...))
			}

			// finished
			d.hideSpinner()
		})
	}
	go call()
}

func (d *MethodCallDialog) showSpinner() {
	d.spinner.Start()
	d.spinner.Display(true)
}

func (d *MethodCallDialog) hideSpinner() {
	d.spinner.Stop()
	d.spinner.Display(false)
}

func (d *MethodCallDialog) Display(display bool) {
	d.display = display
}

func (d *MethodCallDialog) IsDisplay() bool {
	return d.display
}

// Focus implements tview.Focus
func (d *MethodCallDialog) Focus(delegate func(p tview.Primitive)) {
	delegate(d.methods)
}

// SetRect implements tview.SetRect
func (d *MethodCallDialog) SetRect(x int, y int, width int, height int) {
	// self
	dialogWidth := width / 2
	dialogHeight := height / 2
	dialogX := x + ((width - dialogWidth) / 2)
	dialogY := y + ((height - dialogHeight) / 2)
	d.Flex.SetRect(dialogX, dialogY, dialogWidth, dialogHeight)

	// spinner
	d.spinner.SetRect(d.result.GetInnerRect())
}

// Draw implements tview.Primitive
func (d *MethodCallDialog) Draw(screen tcell.Screen) {
	// self
	if d.display {
		d.Flex.Draw(screen)
	}

	// spinner
	d.spinner.Draw(screen)
}
