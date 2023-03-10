package view

import (
	"github.com/dyng/ramen/internal/common"
	"github.com/dyng/ramen/internal/service"
	"github.com/dyng/ramen/internal/view/util"
	"github.com/ethereum/go-ethereum/log"
	"github.com/rivo/tview"
)

type Home struct {
	*tview.Flex
	app *App

	transactionList *TransactionList
}

func NewHome(app *App) *Home {
	home := &Home{
		app: app,
	}

	// setup layout
	home.initLayout()

	// subscribe to new blocks
	app.eventBus.Subscribe(service.TopicNewBlock, home.onNewBlock)

	return home
}

func (h *Home) initLayout() {
	s := h.app.config.Style()

	// Transactions
	transactions := NewTransactionList(h.app, false)
	transactions.SetBorderColor(s.BorderColor)
	transactions.SetTitleColor(s.TitleColor)
	h.transactionList = transactions

	// Root
	flex := tview.NewFlex()
	flex.AddItem(transactions, 0, 1, true)
	h.Flex = flex
}

// KeyMaps implements bodyPage
func (h *Home) KeyMaps() util.KeyMaps {
	return h.transactionList.KeyMaps()
}

func (h *Home) onNewBlock(block *common.Block) {
	txns, err := h.app.service.GetTransactionsByBlock(block)
	if err != nil {
		log.Error("cannot extract transactions from block", "blockHash", block.Hash(), "error", err)
		return
	}

	h.app.QueueUpdateDraw(func() {
		h.transactionList.PrependTransactions(txns)
	})
}
