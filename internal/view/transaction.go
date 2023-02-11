package view

import (
	"fmt"
	"strings"

	"github.com/dyng/ramen/internal/common"
	"github.com/dyng/ramen/internal/common/conv"
	"github.com/dyng/ramen/internal/view/format"
	"github.com/dyng/ramen/internal/view/style"
	"github.com/dyng/ramen/internal/view/util"
	"github.com/ethereum/go-ethereum/log"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	// width of first column
	indentation = len("BlockNumber")
)

type TransactionDetail struct {
	*tview.Flex
	app *App

	transaction common.Transaction
	hash        *util.Section
	blockNumber *util.Section
	timestamp   *util.Section
	from        *util.Section
	to          *util.Section
	value       *util.Section
	data        *util.Section
	calldata    *CallData
}

func NewTransactionDetail(app *App) *TransactionDetail {
	td := &TransactionDetail{
		Flex: tview.NewFlex(),
		app:  app,
	}

	// setup layout
	td.initLayout()

	// setup keymap
	td.initKeymap()

	return td
}

func (t *TransactionDetail) initLayout() {
	s := t.app.config.Style()

	t.SetDirection(tview.FlexRow)
	t.SetBorder(true)
	t.SetTitle(style.BoldPadding("Transaction Detail"))
	t.SetTitleColor(s.TitleColor)
	t.SetBorderColor(s.BorderColor)

	table := tview.NewTable()

	t.hash = util.NewSectionWithStyle("Hash", util.EmptyValue, s)
	t.hash.AddToTable(table, 0, 0)

	t.blockNumber = util.NewSectionWithStyle("BlockNumber", util.EmptyValue, s)
	t.blockNumber.AddToTable(table, 1, 0)

	t.timestamp = util.NewSectionWithStyle("Timestamp", util.EmptyValue, s)
	t.timestamp.AddToTable(table, 2, 0)

	t.from = util.NewSectionWithStyle("From", util.EmptyValue, s)
	t.from.AddToTable(table, 3, 0)

	t.to = util.NewSectionWithStyle("To", util.EmptyValue, s)
	t.to.AddToTable(table, 4, 0)

	t.value = util.NewSectionWithStyle("Value", util.EmptyValue, s)
	t.value.AddToTable(table, 5, 0)

	t.data = util.NewSectionWithStyle("Data", util.EmptyValue, s)
	t.data.AddToTable(table, 6, 0)

	t.calldata = NewCalldata(t.app)

	// add to layout
	t.AddItem(table, 7, 0, false)
	t.AddItem(t.calldata, 0, 1, false)
}

func (t *TransactionDetail) initKeymap() {
	InitKeymap(t, t.app)
}

func (t *TransactionDetail) KeyMaps() util.KeyMaps {
	keymaps := make(util.KeyMaps, 0)

	// KeyF: jump to sender's account page
	keymaps = append(keymaps, util.KeyMap{
		Key:         util.KeyF,
		Shortcut:    "f",
		Description: "To Sender",
		Handler: func(*tcell.EventKey) {
			t.ViewSender()
		},
	})
	// KeyT: jump to receiver's account page
	keymaps = append(keymaps, util.KeyMap{
		Key:         util.KeyT,
		Shortcut:    "t",
		Description: "To Receiver",
		Handler: func(*tcell.EventKey) {
			t.ViewReceiver()
		},
	})

	return keymaps
}

func (t *TransactionDetail) SetTransaction(transaction common.Transaction) {
	t.transaction = transaction
	t.refresh()
}

func (t *TransactionDetail) ViewSender() {
	log.Debug("View transaction sender", "transaction", t.transaction.Hash())
	t.viewAccount(t.from.GetText())
}

func (t *TransactionDetail) ViewReceiver() {
	log.Debug("View transaction receiver", "transaction", t.transaction.Hash())
	t.viewAccount(t.to.GetText())
}

func (t *TransactionDetail) refresh() {
	txn := t.transaction
	t.hash.SetText(txn.Hash().Hex())
	t.blockNumber.SetText(txn.BlockNumber().String())
	t.timestamp.SetText(format.ToDatetime(txn.Timestamp()))
	t.from.SetText(txn.From().Hex())
	t.to.SetText(format.NormalizeReceiverAddress(txn.To()))
	t.value.SetText(fmt.Sprintf("%s (%g Ether)", txn.Value(), conv.ToEther(txn.Value())))
	t.data.SetText(format.BytesToString(txn.Data(), 64))
	t.calldata.LoadAsync(t.transaction.To(), t.transaction.Data())
}

func (t *TransactionDetail) viewAccount(address string) {
	account, err := t.app.service.GetAccount(address)
	if err != nil {
		log.Error("Failed to fetch account of given address", "address", address, "error", err)
		t.app.root.NotifyError(format.FineErrorMessage(
			"Failed to fetch account of address %s", address, err))
	} else {
		t.app.root.ShowAccountPage(account)
	}
}

type CallData struct {
	*tview.Table
	app     *App
	spinner *util.Spinner
}

func NewCalldata(app *App) *CallData {
	c := &CallData{
		Table:   tview.NewTable(),
		app:     app,
		spinner: util.NewSpinner(app.Application),
	}
	c.alignFirstColumn()
	return c
}

func (c *CallData) Clear() {
	c.Table.Clear()
	c.alignFirstColumn()
}

func (c *CallData) LoadAsync(address *common.Address, data []byte) {
	// clear previous data
	c.Clear()

	if len(data) > 0 && address != nil {
		// show spinner
		c.spinner.StartAndShow()

		go func() {
			// populate cache
			_, err := c.app.service.GetContract(*address)
			if err != nil {
				log.Error("Failed to fetch contract", "address", *address, "error", err)
				c.spinner.StopAndHide()
			} else {
				c.app.QueueUpdateDraw(func() {
					hasABI := c.parseData(*address, data)
					if !hasABI {
						c.warnNoABI()
					}
					c.spinner.StopAndHide()
				})
			}
		}()
	}
}

func (c *CallData) parseData(address common.Address, data []byte) bool {
	s := c.app.config.Style()

	if len(data) == 0 {
		return false
	}

	contract, err := c.app.service.GetContract(address)
	if err != nil {
		log.Error("Failed to fetch contract", "address", address, "error", err)
		return false
	}

	if !contract.HasABI() {
		return false
	}

	method, args, err := contract.ParseCalldata(data)
	if err != nil {
		log.Error("Failed to parse calldata", "address", address, "error", err)
		return false
	}

	// set method name
	c.SetCell(0, 1, tview.NewTableCell("[dodgerblue::b]function[-:-:-]"))
	c.SetCell(0, 2, tview.NewTableCell(method.Name).SetAttributes(tcell.AttrBold))

	// set arguments
	for i, argVal := range args {
		arg := method.Inputs[i]

		valStr, err := conv.PackArgument(arg.Type, argVal)
		if err != nil {
			log.Error("Failed to pack argument", "value", argVal, "type", arg.Type, "error", err)
			valStr = "ERROR"
		}

		c.SetCell(i+1, 0, tview.NewTableCell(""))
		c.SetCell(i+1, 1, tview.NewTableCell(arg.Name).SetTextColor(s.SectionColor2))
		c.SetCell(i+1, 2, tview.NewTableCell(valStr))
	}

	return true
}

func (c *CallData) warnNoABI() {
	c.SetCell(0, 1, tview.NewTableCell("[crimson]cannot decode calldata as ABI is unavailable[-]"))
}

func (c *CallData) setSpinnerRect() {
	x, y, _, _ := c.GetInnerRect()
	c.spinner.SetRect(x+indentation+1, y, 0, 0)
}

func (c *CallData) alignFirstColumn() {
	c.SetCell(0, 0, tview.NewTableCell(strings.Repeat(" ", indentation)))
}

// SetRect implements tview.SetRect
func (c *CallData) SetRect(x int, y int, width int, height int) {
	c.Table.SetRect(x, y, width, height)
	c.setSpinnerRect()
}

// Draw implements tview.Primitive
func (c *CallData) Draw(screen tcell.Screen) {
	c.Table.Draw(screen)
	c.spinner.Draw(screen)
}
