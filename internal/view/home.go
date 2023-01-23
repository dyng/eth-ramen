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
	transactions := NewTransactionList(h.app)
	transactions.SetBorderColor(s.PrimaryBorderColor)
	transactions.SetTitleColor(s.PrimaryTitleColor)
	h.transactionList = transactions

	// Root
	flex := tview.NewFlex()
	flex.AddItem(transactions, 0, 1, true)
	h.Flex = flex
}

// KeyMaps implements bodyPage
func (h *Home) KeyMaps() util.KeyMaps {
	keymaps := make(util.KeyMaps, 0)
	return keymaps
}

func (h *Home) onNewBlock(block *common.Block) {
	h.app.QueueUpdateDraw(func() {
		txns, err := h.app.service.GetTransactionsByBlock(block)
		if err != nil {
			log.Error("cannot extract transactions from block", "blockHash", block.Hash(), "error", err)
			return
		}
		h.transactionList.SetTransactions(txns)
	})
}
