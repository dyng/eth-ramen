package view

import (
	"fmt"
	"math/big"

	"github.com/dyng/ramen/internal/common/conv"
	"github.com/dyng/ramen/internal/service"
	"github.com/dyng/ramen/internal/view/format"
	"github.com/dyng/ramen/internal/view/style"
	"github.com/dyng/ramen/internal/view/util"
	gcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	// transferDialogMinHeight is the minimum height of the transfer dialog.
	transferDialogMinHeight = 10
	// transferDialogMinWidth is the minimum width of the transfer dialog.
	transferDialogMinWidth  = 50
)

type TransferDialog struct {
	*tview.Form
	app       *App
	display   bool
	lastFocus tview.Primitive

	sender *service.Signer
	info   *SenderFormItem
	to     *tview.InputField
	amount *tview.InputField
}

func NewTransferDialog(app *App) *TransferDialog {
	d := &TransferDialog{
		app:     app,
		display: false,
	}

	// setup layout
	d.initLayout()

	// setup keymap
	d.initKeymap()

	return d
}

func (d *TransferDialog) initLayout() {
	s := d.app.config.Style()

	// sender info
	info := NewSenderFormItem(d.app)
	d.info = info

	// form
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetBorderColor(s.DialogBorderColor)
	form.SetTitle(style.BoldPadding("Transfer"))
	form.SetLabelColor(s.InputFieldLableColor)
	form.SetFieldBackgroundColor(s.InputFieldBgColor)
	form.SetButtonsAlign(tview.AlignRight)
	form.SetButtonBackgroundColor(s.ButtonBgColor)
	form.AddFormItem(info)
	form.AddInputField("To", "", 999, nil, nil)
	form.AddInputField("Amount", "", 999, nil, nil)
	form.AddButton("Transfer", d.doTransfer)
	d.to = form.GetFormItemByLabel("To").(*tview.InputField)
	d.amount = form.GetFormItemByLabel("Amount").(*tview.InputField)
	d.Form = form
}

func (d *TransferDialog) initKeymap() {
	InitKeymap(d, d.app)
}

// KeyMaps implements KeymapPrimitive
func (d *TransferDialog) KeyMaps() util.KeyMaps {
	keymaps := make(util.KeyMaps, 0)
	keymaps = append(keymaps, util.NewSimpleKey(tcell.KeyEsc, d.Hide))
	return keymaps
}

func (d *TransferDialog) SetSender(account *service.Signer) {
	d.sender = account
	d.refresh()
}

func (d *TransferDialog) refresh() {
	// refresh sender's information (e.g. balance)
	d.info.SetSender(d.sender)
}

// doTransfer is core method that do the whole things
func (d *TransferDialog) doTransfer() {
	i, ok := new(big.Float).SetString(d.amount.GetText())
	if !ok {
		d.app.root.NotifyError(fmt.Sprintf("Cannot parse amount value %s", d.amount.GetText()))
		return
	}

	// close dialog
	d.Hide()

	amount := conv.FromEther(i)
	toAddr := gcommon.HexToAddress(d.to.GetText())
	log.Info("Transfer ethers to another account", "from", d.sender.GetAddress(), "to", toAddr, "amount", amount)

	hash, err := d.sender.TransferTo(toAddr, amount)
	if err != nil {
		d.app.root.NotifyError(format.FineErrorMessage("Failed to complete transfer", err))
	} else {
		d.app.root.NotifyInfo(fmt.Sprintf("Transaction has been submitted.\n\nTxnHash: %s", hash))
	}
}

func (d *TransferDialog) Show() {
	if !d.display {
		// save last focused element
		d.lastFocus = d.app.GetFocus()

		d.Display(true)
		d.app.SetFocus(d)
	}
}

func (d *TransferDialog) Hide() {
	if d.display {
		d.Display(false)
		d.app.SetFocus(d.lastFocus)
	}
}

func (d *TransferDialog) ClearAndRefresh() {
	// clear
	d.to.SetText("")
	d.amount.SetText("")

	// refresh
	d.refresh()
}

func (d *TransferDialog) Display(display bool) {
	d.display = display
}

func (d *TransferDialog) IsDisplay() bool {
	return d.display
}

// Draw implements tview.Primitive
func (d *TransferDialog) Draw(screen tcell.Screen) {
	if d.display {
		d.Form.Draw(screen)
	}
}

func (d *TransferDialog) SetCentral(x int, y int, width int, height int) {
	dialogWidth := width - width/2
	dialogHeight := style.AvatarSize + 12
	if dialogHeight < transferDialogMinHeight {
		dialogHeight = transferDialogMinHeight
	}
	if dialogWidth < transferDialogMinWidth {
		dialogWidth = transferDialogMinWidth
	}
	dialogX := x + ((width - dialogWidth) / 2)
	dialogY := y + ((height - dialogHeight) / 2)
	d.Form.SetRect(dialogX, dialogY, dialogWidth, dialogHeight)
}

type SenderFormItem struct {
	*tview.Flex
	app *App

	label     *tview.Table
	lableCell *tview.TableCell
	field     *tview.Flex
	avatar    *util.Avatar
	address   *util.Section
	balance   *util.Section
}

func NewSenderFormItem(app *App) *SenderFormItem {
	fi := &SenderFormItem{
		app:    app,
		avatar: util.NewAvatar(style.AvatarSize),
	}
	s := app.config.Style()
	table := tview.NewTable()

	// address
	address := util.NewSectionWithStyle("Address", util.EmptyValue, s)
	address.AddToTable(table, 0, 0)
	fi.address = address

	// balance
	balance := util.NewSectionWithStyle("Balance", util.EmptyValue, s)
	balance.AddToTable(table, 1, 0)
	fi.balance = balance

	// field
	field := tview.NewFlex().SetDirection(tview.FlexRow)
	field.AddItem(fi.avatar, style.AvatarSize, 0, false)
	field.AddItem(table, 2, 0, false)
	fi.field = field

	// label
	label := tview.NewTable()
	cell := tview.NewTableCell(fi.GetLabel())
	label.SetCell((fi.GetFieldHeight()-1)/2, 0, cell)
	fi.label = label
	fi.lableCell = cell

	// flex
	flex := tview.NewFlex().SetDirection(tview.FlexColumn)
	flex.AddItem(label, 1, 0, false)
	flex.AddItem(field, 0, 1, false)
	fi.Flex = flex

	return fi
}

func (s *SenderFormItem) SetSender(account *service.Signer) {
	addr := account.GetAddress()

	// avatar
	s.avatar.SetAddress(addr)

	// address
	s.address.SetText(addr.Hex())

	// balance
	bal := account.GetBalance()
	s.balance.SetText(conv.ToEther(bal).String())
}

// Focus implements tview.Primitive
func (s *SenderFormItem) Focus(delegate func(p tview.Primitive)) {
	delegate(s.app.root.transfer.GetFormItemByLabel("To"))
}

// GetFieldHeight implements tview.FormItem
func (s *SenderFormItem) GetFieldHeight() int {
	return style.AvatarSize + 2
}

// GetFieldWidth implements tview.FormItem
func (s *SenderFormItem) GetFieldWidth() int {
	return 999
}

// GetLabel implements tview.FormItem
func (s *SenderFormItem) GetLabel() string {
	return "From"
}

// SetFinishedFunc implements tview.FormItem
func (s *SenderFormItem) SetFinishedFunc(handler func(key tcell.Key)) tview.FormItem {
	return s
}

// SetFormAttributes implements tview.FormItem
func (s *SenderFormItem) SetFormAttributes(labelWidth int, labelColor tcell.Color, bgColor tcell.Color, fieldTextColor tcell.Color, fieldBgColor tcell.Color) tview.FormItem {
	s.lableCell.SetTextColor(labelColor)
	s.lableCell.SetBackgroundColor(bgColor)
	s.Flex.ResizeItem(s.label, labelWidth, 0)
	return s
}
