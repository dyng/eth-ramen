package view

import (
	"strings"

	"github.com/dyng/ramen/internal/common"
	"github.com/dyng/ramen/internal/common/conv"
	"github.com/dyng/ramen/internal/view/format"
	"github.com/rivo/tview"
)

type TransactionList struct {
	*tview.Table
	app *App

	txns common.Transactions
}

func NewTransactionList(app *App) *TransactionList {
	t := &TransactionList{
		Table: tview.NewTable(),
		app:   app,
		txns:  []common.Transaction{},
	}

	// setup layout
	t.initLayout()

	return t
}

func (t *TransactionList) initLayout() {
	t.SetBorder(true)
	t.SetTitle(" Transactions ")

	headers := []string{"hash", "block", "from", "to", "value", "datetime"}
	for i, header := range headers {
		t.SetCell(0, i,
			tview.NewTableCell(strings.ToUpper(header)).
				SetExpansion(1).
				SetAlign(tview.AlignLeft).
				SetSelectable(false))
	}

	t.SetSelectable(true, false)
	t.SetFixed(1, 1)
	t.SetSelectedFunc(t.handleSelected)
}

func (t *TransactionList) SetTransactions(txns common.Transactions) {
	t.txns = txns
	t.refresh()
}

func (t *TransactionList) Clear() {
	for i := t.GetRowCount() - 1; i > 0; i-- {
		t.RemoveRow(i)
	}
}

func (t *TransactionList) refresh() {
	// clear previous content at first
	t.Clear()

	for i := 0; i < len(t.txns); i++ {
		tx := t.txns[i]
		row := i + 1
		t.SetCell(row, 0, tview.NewTableCell(format.TruncateText(tx.Hash().Hex(), 8)))
		t.SetCell(row, 1, tview.NewTableCell(tx.BlockNumber().String()))
		t.SetCell(row, 2, tview.NewTableCell(format.TruncateText(tx.From().Hex(), 20)))
		t.SetCell(row, 3, tview.NewTableCell(format.TruncateText(
			format.NormalizeReceiverAddress(tx.To()), 20)))
		t.SetCell(row, 4, tview.NewTableCell(conv.ToEther(tx.Value()).String()))
		t.SetCell(row, 5, tview.NewTableCell(format.ToDatetime(tx.Timestamp())))
	}
}

func (t *TransactionList) handleSelected(row int, column int) {
	if row > 0 && row <= len(t.txns) {
		txn := t.txns[row-1]
		t.app.root.ShowTransactionPage(txn)
	}
}
