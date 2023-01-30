package view

import (
	"fmt"

	"github.com/dyng/ramen/internal/common"
	"github.com/dyng/ramen/internal/common/conv"
	"github.com/dyng/ramen/internal/view/format"
	"github.com/dyng/ramen/internal/view/style"
	"github.com/dyng/ramen/internal/view/util"
	"github.com/ethereum/go-ethereum/log"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TransactionDetail struct {
	*tview.Table
	app *App

	transaction common.Transaction
	hash        *util.Section
	blockNumber *util.Section
	timestamp   *util.Section
	from        *util.Section
	to          *util.Section
	value       *util.Section
}

func NewTransactionDetail(app *App) *TransactionDetail {
	td := &TransactionDetail{
		Table: tview.NewTable(),
		app:   app,
	}

	// setup layout
	td.initLayout()

	// setup keymap
	td.initKeymap()

	return td
}

func (t *TransactionDetail) initLayout() {
	s := t.app.config.Style()

	t.SetBorder(true)
	t.SetTitle(style.BoldPadding("Transaction Detail"))
	t.SetTitleColor(s.TitleColor)
	t.SetBorderColor(s.BorderColor)

	t.hash = util.NewSectionWithColor("Hash", s.SectionColor, util.EmptyValue, s.FgColor)
	t.hash.AddToTable(t.Table, 0, 0)

	t.blockNumber = util.NewSectionWithColor("BlockNumber", s.SectionColor, util.EmptyValue, s.FgColor)
	t.blockNumber.AddToTable(t.Table, 1, 0)

	t.timestamp = util.NewSectionWithColor("Timestamp", s.SectionColor, util.EmptyValue, s.FgColor)
	t.timestamp.AddToTable(t.Table, 2, 0)

	t.from = util.NewSectionWithColor("From", s.SectionColor, util.EmptyValue, s.FgColor)
	t.from.AddToTable(t.Table, 3, 0)

	t.to = util.NewSectionWithColor("To", s.SectionColor, util.EmptyValue, s.FgColor)
	t.to.AddToTable(t.Table, 4, 0)

	t.value = util.NewSectionWithColor("Value", s.SectionColor, util.EmptyValue, s.FgColor)
	t.value.AddToTable(t.Table, 5, 0)
}

func (t *TransactionDetail) initKeymap() {
	InitKeymap(t, t.app)
}

func (t *TransactionDetail) KeyMaps() util.KeyMaps {
	keymaps := make(util.KeyMaps, 0)

	// KeyF: jump to sender's account page
	keymaps = append(keymaps, util.KeyMap{
		Key:         util.KeyF,
		Shortcut:    "F",
		Description: "View Sender",
		Handler: func(*tcell.EventKey) {
			t.viewAccount(t.from.GetText())
		},
	})
	// KeyT: jump to receiver's account page
	keymaps = append(keymaps, util.KeyMap{
		Key:         util.KeyT,
		Shortcut:    "T",
		Description: "View Receiver",
		Handler: func(*tcell.EventKey) {
			t.viewAccount(t.to.GetText())
		},
	})

	return keymaps
}

func (t *TransactionDetail) SetTransaction(transaction common.Transaction) {
	t.transaction = transaction
	t.refresh()
}

func (t *TransactionDetail) refresh() {
	txn := t.transaction
	t.hash.SetText(txn.Hash().Hex())
	t.blockNumber.SetText(txn.BlockNumber().String())
	t.timestamp.SetText(format.ToDatetime(txn.Timestamp()))
	t.from.SetText(txn.From().Hex())
	t.to.SetText(format.NormalizeReceiverAddress(txn.To()))
	t.value.SetText(fmt.Sprintf("%s (%g Ether)", txn.Value(), conv.ToEther(txn.Value())))
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
