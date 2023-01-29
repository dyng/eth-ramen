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
	app    *App
	loader *util.Loader

	showInOut bool
	base      *common.Address
	txns      common.Transactions
}

func NewTransactionList(app *App, showInOut bool) *TransactionList {
	t := &TransactionList{
		Table:     tview.NewTable(),
		app:       app,
		loader:    util.NewLoader(app.Application),
		showInOut: showInOut,
		txns:      []common.Transaction{},
	}

	// setup layout
	t.initLayout()

	return t
}

func (t *TransactionList) initLayout() {
	s := t.app.config.Style()

	t.SetBorder(true)
	t.SetTitle(style.BoldPadding("Transactions"))

	// table
	var headers []string
	if  t.showInOut {
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

func (t *TransactionList) SetBaseAccount(account *common.Address) {
	t.base = account
}

func (t *TransactionList) FilterAndPrependTransactions(txns common.Transactions) {
	if t.base == nil {
		t.PrependTransactions(txns)
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

func (t *TransactionList) PrependTransactions(txns common.Transactions) {
	prepended := append(txns, t.txns...)
	t.SetTransactions(prepended)
}

func (t *TransactionList) SetTransactions(txns common.Transactions) {
	if len(txns) > TransactionListLimit {
		txns = txns[:TransactionListLimit]
	}
	t.txns = txns
	t.refresh()
}

func (t *TransactionList) LoadAsync(loader func() (common.Transactions, error)) {
	// clear current content
	t.Clear()

	// display loader
	t.loader.Start()
	t.loader.Display(true)

	load := func() {
		txns, err := loader()
		t.app.QueueUpdateDraw(func() {
			if err == nil {
				t.loader.Stop()
				t.loader.Display(false)
				if txns != nil {
					t.SetTransactions(txns)
				}
			} else {
				log.Error("Failed to load transactions", "error", err)
				t.app.root.NotifyError(format.FineErrorMessage("Error occurs when loading transactions.", err))
			}
		})
	}

	go load()
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

func (t *TransactionList) handleSelected(row int, column int) {
	if row > 0 && row <= len(t.txns) {
		txn := t.txns[row-1]
		t.app.root.ShowTransactionPage(txn)
	}
}

// SetRect implements tview.SetRect
func (t *TransactionList) SetRect(x, y, width, height int) {
	t.Table.SetRect(x, y, width, height)
	t.loader.SetCentral(x, y, width, height)
}

// Draw implements tview.Draw
func (t *TransactionList) Draw(screen tcell.Screen) {
	t.Table.Draw(screen)
	t.loader.Draw(screen)
}
