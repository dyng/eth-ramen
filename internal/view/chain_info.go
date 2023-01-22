package view

import (
	"fmt"

	"github.com/dyng/ramen/internal/common"
	"github.com/dyng/ramen/internal/service"
	"github.com/dyng/ramen/internal/view/util"
	"github.com/rivo/tview"
)

type ChainInfo struct {
	*tview.Table
	app *App

	network  *util.Section
	height   *util.Section
	gasPrice *util.Section
}

func NewChainInfo(app *App) *ChainInfo {
	chainInfo := &ChainInfo{
		Table: tview.NewTable(),
		app: app,
	}

	chainInfo.network = util.NewSection("Network", util.NAValue)
	chainInfo.network.AddToTable(chainInfo.Table, 0, 0)

	chainInfo.height = util.NewSection("Block Height", util.NAValue)
	chainInfo.height.AddToTable(chainInfo.Table, 1, 0)

	chainInfo.gasPrice = util.NewSection("Gas Price", util.NAValue)
	chainInfo.gasPrice.AddToTable(chainInfo.Table, 2, 0)

	chainInfo.app.eventBus.Subscribe(service.TopicNewBlock, chainInfo.onNewBlock)

	return chainInfo
}

func (ci *ChainInfo) SetNetwork(network string) {
	ci.network.SetText(network)
}

func (ci *ChainInfo) SetHeight(height uint64) {
	ci.height.SetText(fmt.Sprint(height))
}

func (ci *ChainInfo) SetGasPrice(gasPrice common.BigInt)  {
	ci.gasPrice.SetText(gasPrice.String())
}

func (ci *ChainInfo) onNewBlock(block *common.Block) {
	ci.app.QueueUpdateDraw(func() {
		ci.SetHeight(block.Number().Uint64())
	})
}
