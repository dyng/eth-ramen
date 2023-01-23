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

	s := app.config.Style()

	chainInfo.network = util.NewSectionWithColor("Network:", s.SectionColor, util.NAValue, s.FgColor)
	chainInfo.network.AddToTable(chainInfo.Table, 0, 0)

	chainInfo.height = util.NewSectionWithColor("Block Height:", s.SectionColor, util.NAValue, s.FgColor)
	chainInfo.height.AddToTable(chainInfo.Table, 1, 0)

	chainInfo.gasPrice = util.NewSectionWithColor("Gas Price:", s.SectionColor, util.NAValue, s.FgColor)
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
