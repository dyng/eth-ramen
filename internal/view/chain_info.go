package view

import (
	"fmt"

	"github.com/dyng/ramen/internal/common"
	"github.com/dyng/ramen/internal/common/conv"
	"github.com/dyng/ramen/internal/service"
	"github.com/dyng/ramen/internal/view/util"
	"github.com/rivo/tview"
	"github.com/shopspring/decimal"
)

type ChainInfo struct {
	*tview.Table
	app *App

	network   *util.Section
	height    *util.Section
	gasPrice  *util.Section
	ethPrice  *util.Section
	prevPrice *decimal.Decimal
}

func NewChainInfo(app *App) *ChainInfo {
	chainInfo := &ChainInfo{
		Table: tview.NewTable(),
		app:   app,
	}

	// setup layout
	chainInfo.initLayout()

	// subscribe for new data
	chainInfo.app.eventBus.Subscribe(service.TopicNewBlock, chainInfo.onNewBlock)
	chainInfo.app.eventBus.Subscribe(service.TopicChainData, chainInfo.onNewChainData)

	return chainInfo
}

func (ci *ChainInfo) initLayout() {
	s := ci.app.config.Style()

	network := util.NewSectionWithColor("Network:", s.SectionColor, util.NAValue, s.FgColor)
	network.AddToTable(ci.Table, 0, 0)
	ci.network = network

	height := util.NewSectionWithColor("Block Height:", s.SectionColor, util.NAValue, s.FgColor)
	height.AddToTable(ci.Table, 1, 0)
	ci.height = height

	gasPrice := util.NewSectionWithColor("Gas Price:", s.SectionColor, util.NAValue, s.FgColor)
	gasPrice.AddToTable(ci.Table, 2, 0)
	ci.gasPrice = gasPrice

	ethPrice := util.NewSectionWithColor("Ether:", s.SectionColor, util.NAValue, s.FgColor)
	ethPrice.AddToTable(ci.Table, 0, 2)
	ci.ethPrice = ethPrice
}

func (ci *ChainInfo) SetNetwork(network string) {
	ci.network.SetText(network)
}

func (ci *ChainInfo) SetHeight(height uint64) {
	ci.height.SetText(fmt.Sprint(height))
}

func (ci *ChainInfo) SetGasPrice(gasPrice common.BigInt) {
	ci.gasPrice.SetText(fmt.Sprintf("%s Gwei", conv.ToGwei(gasPrice)))
}

func (ci *ChainInfo) SetEthPrice(price decimal.Decimal) {
	if ci.prevPrice == nil {
		ci.ethPrice.SetText(fmt.Sprintf("$%s", price))
	} else {
		c := ci.prevPrice.Cmp(price)
		if c == 0 {
			// if price does not change, don't change anything
			return
		}

		if c < 0 {
			ci.ethPrice.SetText(fmt.Sprintf("[lightgreen]$%s ▲[-]", price))
		} else {
			ci.ethPrice.SetText(fmt.Sprintf("[crimson]$%s ▼[-]", price))
		}
	}

	ci.prevPrice = &price
}

func (ci *ChainInfo) onNewBlock(block *common.Block) {
	ci.app.QueueUpdateDraw(func() {
		ci.SetHeight(block.Number().Uint64())
	})
}

func (ci *ChainInfo) onNewChainData(data *service.ChainData) {
	ci.app.QueueUpdateDraw(func() {
		if data.Price != nil {
			ci.SetEthPrice(*data.Price)
		}
		if data.GasPrice != nil {
			ci.SetGasPrice(data.GasPrice)
		}
	})
}
