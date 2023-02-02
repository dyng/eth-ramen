package view

import (
	"fmt"
	"sort"

	"github.com/dyng/ramen/internal/common/conv"
	"github.com/dyng/ramen/internal/service"
	"github.com/dyng/ramen/internal/view/format"
	"github.com/dyng/ramen/internal/view/style"
	"github.com/dyng/ramen/internal/view/util"
	"github.com/ethereum/go-ethereum/accounts/abi"
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
		app:     app,
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
	methods.SetBorderColor(s.DialogBorderColor)
	methods.SetTitle(style.Padding("Method"))
	methods.SetSelectable(true, false)
	methods.SetSelectionChangedFunc(func(row, column int) {
		d.showArguments()
	})
	methods.SetSelectedFunc(func(row, column int) {
		if d.methodHasNoArg() {
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
	args.SetBorderColor(s.DialogBorderColor)
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
	d.methods.Clear()
	d.args.Clear(true)
	d.result.Clear()
	if d.contract.HasABI() {
		d.showMethodList()
	}
}

func (d *MethodCallDialog) showMethodList() {
	d.methods.Clear()

	s := d.app.config.Style()
	row := 0
	for _, method := range d.sortedMethods() {
		if method.IsConstant() {
			color := s.FgColor
			d.methods.SetCell(row, 0, tview.NewTableCell(method.Name).SetTextColor(color).SetExpansion(1))
			d.methods.SetCell(row, 1, tview.NewTableCell(" ").SetTextColor(color))
		} else {
			color := tcell.ColorDarkRed
			d.methods.SetCell(row, 0, tview.NewTableCell(method.Name).SetExpansion(1).SetBackgroundColor(color))
			d.methods.SetCell(row, 1, tview.NewTableCell("⚠").SetBackgroundColor(color))
		}
		row++
	}
}

// sortedMethods returns a list of method sorted by name and purity
func (d *MethodCallDialog) sortedMethods() []abi.Method {
	methods := d.contract.GetABI().Methods
	sorted := make([]abi.Method, 0)
	for _, m := range methods {
		sorted = append(sorted, m)
	}

	sort.Slice(sorted, func(i, j int) bool {
		m1 := sorted[i] 
		m2 := sorted[j]

		if m1.IsConstant() {
			if m2.IsConstant() {
				return m1.Name < m2.Name
			} else {
				return true
			}
		} else {
			if m2.IsConstant() {
				return false
			} else {
				return m1.Name < m2.Name
			}
		}
	})

	return sorted
}

func (d *MethodCallDialog) showArguments() {
	d.args.Clear(true)

	row, _ := d.methods.GetSelection()
	methodName := d.methods.GetCell(row, 0).Text
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

func (d *MethodCallDialog) methodHasNoArg() bool {
	methodName := d.methods.GetCell(d.methods.GetSelection()).Text
	method := d.contract.GetABI().Methods[methodName]
	return len(method.Inputs) == 0
}

func (d *MethodCallDialog) callMethod() {
	// start calling
	d.spinner.StartAndShow()

	methodName := d.methods.GetCell(d.methods.GetSelection()).Text
	method := d.contract.GetABI().Methods[methodName]

	// unpack arguments
	args := make([]any, 0)
	for i := 0; i < d.args.GetFormItemCount(); i++ {
		item := d.args.GetFormItem(i).(*tview.InputField)
		arg := method.Inputs[i]
		val, err := conv.UnpackArgument(arg.Type, item.GetText())
		if err == nil {
			args = append(args, val)
		} else {
			log.Error("Cannot unpack argument", "argument", arg, "input", item.GetText(), "error", err)
			d.app.root.NotifyError(format.FineErrorMessage(
				"Input type for argument '%s' is incorrect, should be '%s'.", arg.Name, arg.Type.String(), err))
			return
		}
	}

	// ensure signer has signed in
	var signer *service.Signer
	if !method.IsConstant() {
		signer = d.app.root.signer.GetSigner()
		if signer == nil {
			d.app.root.NotifyError(format.FineErrorMessage(""))
			return
		}
	}

	go func() {
		var res []any
		var err error

		if method.IsConstant() {
			res, err = d.contract.Call(methodName, args...)
		} else {
			// FIXME: waiting for transaction executed
			hash, e := d.contract.Transact(signer, methodName, args...)
			res = []any{fmt.Sprintf("Submitted! TxnHash: %s", hash)}
			err = e
		}

		d.app.QueueUpdateDraw(func() {
			if err != nil {
				log.Error("Method call is failed", "name", methodName, "args", args, "error", err)
				d.app.root.NotifyError(format.FineErrorMessage("Cannot call contract method '%s'.", methodName, err))
			} else {
				d.result.SetText(fmt.Sprint(res...))
			}

			// calling finished
			d.spinner.StopAndHide()
		})
	}()
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

func (d *MethodCallDialog) SetCentral(x int, y int, width int, height int) {
	// self
	dialogWidth := width / 2
	dialogHeight := height / 2
	dialogX := x + ((width - dialogWidth) / 2)
	dialogY := y + ((height - dialogHeight) / 2)
	d.Flex.SetRect(dialogX, dialogY, dialogWidth, dialogHeight)

	// spinner
	d.spinner.SetCentral(d.result.GetInnerRect())
}

// Draw implements tview.Primitive
func (d *MethodCallDialog) Draw(screen tcell.Screen) {
	if d.display {
		d.Flex.Draw(screen)
	}
	d.spinner.Draw(screen)
}
