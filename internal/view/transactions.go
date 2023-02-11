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

const (
	// TransactionListLimit is the length limit of transaction list
	TransactionListLimit = 1000
)

type TransactionList struct {
	*tview.Table
	app     *App
	txnPrev *TxnPreviewDialog
	loader  *util.Loader

	showInOut bool
	base      *common.Address
	txns      common.Transactions
}

func NewTransactionList(app *App, showInOut bool) *TransactionList {
	t := &TransactionList{
		Table:     tview.NewTable(),
		app:       app,
		txnPrev:   NewTxnPreviewDialog(app),
		loader:    util.NewLoader(app.Application),
		showInOut: showInOut,
		txns:      []common.Transaction{},
	}

	// setup layout
	t.initLayout()

	// setup keymap
	t.initKeymap()

	return t
}

func (t *TransactionList) initLayout() {
	s := t.app.config.Style()

	t.SetBorder(true)
	t.SetTitle(style.BoldPadding("Transactions"))

	// table
	var headers []string
	if t.showInOut {
		headers = []string{"hash", "block", "from", "to", "", "value", "datetime"}
	} else {
		headers = []string{"hash", "block", "from", "to", "value", "datetime"}
	}
	for i, header := range headers {
		t.SetCell(0, i,
			tview.NewTableCell(strings.ToUpper(header)).
				SetExpansion(1).
				SetAlign(tview.AlignLeft).
				SetStyle(s.TableHeaderStyle).
				SetSelectable(false))
	}
	t.SetSelectable(true, false)
	t.SetFixed(1, 1)
	t.SetSelectedFunc(t.handleSelected)

	// loader
	t.loader.SetTitleColor(s.PrgBarTitleColor)
	t.loader.SetBorderColor(s.PrgBarBorderColor)
	t.loader.SetCellColor(s.PrgBarCellColor)
}

func (t *TransactionList) initKeymap() {
	InitKeymap(t, t.app)
}

func (t *TransactionList) KeyMaps() util.KeyMaps {
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

// SetBaseAccount sets the base account to determine whether a transaction is
// inflow or outflow
func (t *TransactionList) SetBaseAccount(account *common.Address) {
	t.base = account
}

// FilterAndPrependTransactions is like PrependTransactions, but filters out
// transactions that are not related to the base account
func (t *TransactionList) FilterAndPrependTransactions(txns common.Transactions) {
	if t.base == nil {
		return
	}

	toAdd := make(common.Transactions, 0)
	for _, tx := range txns {
		if tx.From().String() == t.base.String() {
			toAdd = append(toAdd, tx)
		}
		if tx.To() != nil && tx.To().String() == t.base.String() {
			toAdd = append(toAdd, tx)
		}
	}
	t.PrependTransactions(toAdd)
}

// PrependTransactions prepends transactions to existing transactions
func (t *TransactionList) PrependTransactions(txns common.Transactions) {
	prepended := append(txns, t.txns...)
	t.SetTransactions(prepended)
}

// SetTransactions sets a transaction list
func (t *TransactionList) SetTransactions(txns common.Transactions) {
	if len(txns) > TransactionListLimit {
		txns = txns[:TransactionListLimit]
	}
	t.txns = txns
	t.refresh()
}

// LoadAsync loads transactions asynchronously
func (t *TransactionList) LoadAsync(loader func() (common.Transactions, error)) {
	// clear current content
	t.Clear()

	// start loading animation
	t.loader.Start()
	t.loader.Display(true)

	go func() {
		txns, err := loader()
		t.app.QueueUpdateDraw(func() {
			// stop loading animation
			t.loader.Stop()
			t.loader.Display(false)

			if err == nil {
				if txns != nil {
					t.SetTransactions(txns)
				}
			} else {
				log.Error("Failed to load transactions", "error", err)
				t.app.root.NotifyError(format.FineErrorMessage("Error occurs when loading transactions.", err))
			}
		})
	}()
}

// ViewSender jumps to the sender's account page
func (t *TransactionList) ViewSender() {
	current := t.selection()
	if current == nil {
		return
	}

	addr := current.From()
	if t.base != nil && addr.String() == t.base.String() {
		return
	}

	t.viewAccount(addr)
}

// ViewReceiver jumps to the receiver's account page
func (t *TransactionList) ViewReceiver() {
	current := t.selection()
	if current == nil {
		return
	}

	addr := current.To()
	if t.base != nil && addr.String() == t.base.String() {
		return
	}

	t.viewAccount(addr)
}

func (t *TransactionList) Clear() {
	for i := t.GetRowCount() - 1; i > 0; i-- {
		t.RemoveRow(i)
	}
}

func (t *TransactionList) refresh() {
	// clear previous content at first
	t.Clear()

	// show transaction count
	t.SetTitle(style.BoldPadding(fmt.Sprintf("Transactions[[coral]%d[-]]", len(t.txns))))

	for i := 0; i < len(t.txns); i++ {
		tx := t.txns[i]
		row := i + 1

		j := 0
		t.SetCell(row, Inc(&j), tview.NewTableCell(format.TruncateText(tx.Hash().Hex(), 8)))
		t.SetCell(row, Inc(&j), tview.NewTableCell(tx.BlockNumber().String()))
		t.SetCell(row, Inc(&j), tview.NewTableCell(format.TruncateText(tx.From().Hex(), 20)))
		t.SetCell(row, Inc(&j), tview.NewTableCell(format.TruncateText(
			format.NormalizeReceiverAddress(tx.To()), 20)))
		if t.showInOut {
			t.SetCell(row, Inc(&j), tview.NewTableCell(StyledTxnDirection(t.base, tx)))
		}
		t.SetCell(row, Inc(&j), tview.NewTableCell(conv.ToEther(tx.Value()).String()))
		t.SetCell(row, Inc(&j), tview.NewTableCell(format.ToDatetime(tx.Timestamp())))
	}
}

// handleSelected shows a preview of selected transaction
func (t *TransactionList) handleSelected(row int, column int) {
	if row > 0 && row <= len(t.txns) {
		txn := t.txns[row-1]
		t.txnPrev.SetTransaction(txn)
		t.txnPrev.Show()
	}
}

func (t *TransactionList) viewAccount(address *common.Address) {
	if address == nil {
		return
	}

	account, err := t.app.service.GetAccount(address.Hex())
	if err != nil {
		log.Error("Failed to fetch account of given address", "address", address.Hex(), "error", err)
		t.app.root.NotifyError(format.FineErrorMessage(
			"Failed to fetch account of address %s", address.Hex(), err))
	} else {
		t.txnPrev.Hide() // hide dialog if it's visible
		t.app.root.ShowAccountPage(account)
	}
}

func (t *TransactionList) selection() common.Transaction {
	row, _ := t.GetSelection()
	if row > 0 && row <= len(t.txns) {
		return t.txns[row-1]
	} else {
		return nil
	}
}

// HasFocus implements tview.Primitive
func (t *TransactionList) HasFocus() bool {
	if t.txnPrev.HasFocus() {
		return true
	}
	return t.Table.HasFocus()
}

// InputHandler implements tview.Primitive
func (t *TransactionList) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if t.txnPrev.HasFocus() {
			if handler := t.txnPrev.InputHandler(); handler != nil {
				handler(event, setFocus)
				return
			}
		}
		if t.Table.HasFocus() {
			if handler := t.Table.InputHandler(); handler != nil {
				handler(event, setFocus)
				return
			}
		}
	}
}

// SetRect implements tview.SetRect
func (t *TransactionList) SetRect(x, y, width, height int) {
	t.Table.SetRect(x, y, width, height)
	t.txnPrev.SetCentral(x, y, width, height)
	t.loader.SetCentral(x, y, width, height)
}

// Draw implements tview.Draw
func (t *TransactionList) Draw(screen tcell.Screen) {
	t.Table.Draw(screen)
	t.txnPrev.Draw(screen)
	t.loader.Draw(screen)
}
