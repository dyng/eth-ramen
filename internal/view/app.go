package view

import (
	"github.com/asaskevich/EventBus"
	"github.com/dyng/ramen/internal/common"
	conf "github.com/dyng/ramen/internal/config"
	serv "github.com/dyng/ramen/internal/service"
	"github.com/ethereum/go-ethereum/log"
	"github.com/rivo/tview"
)

type App struct {
	*tview.Application
	root *Root

	service  *serv.Service
	config   *conf.Config
	eventBus EventBus.Bus
	syncer   *serv.Syncer
}

func NewApp(config *conf.Config) *App {
	log.Info("Start application with configurations", "config", config)

	app := &App{
		Application: tview.NewApplication(),
		config:      config,
		eventBus:    EventBus.New(),
		service:     serv.NewService(config),
	}

	// syncer
	syncer := serv.NewSyncer(app.service, app.eventBus)
	app.syncer = syncer

	// root
	root := NewRoot(app)
	app.root = root
	app.SetRoot(root, true)

	return app
}

func (a *App) Start() error {
	// first synchronization at startup
	err := a.firstSync()
	if err != nil {
		log.Crit("Failed to synchronize chain info", "error", err)
	}

	// show homepage
	a.root.ShowHomePage()

	// start application
	log.Info("Application is running")
	return a.Run()
}

// firstSync synchronize latest blockchain informations and populate data to widgets
func (a *App) firstSync() error {
	// update network
	network := a.service.GetNetwork()
	a.root.chainInfo.SetNetwork(StyledNetworkName(network))

	// update block height
	go func() {
		height, err := a.service.GetBlockHeight()
		if err != nil {
			log.Error("Failed to fetch block height", "error", err)
		}

		price, err := a.service.GetEthPrice()
		if err != nil {
			log.Error("Failed to fetch ether's price", "error", err)
		}

		gasPrice, err := a.service.GetGasPrice()
		if err != nil {
			log.Error("Failed to fetch gas price", "error", err)
		}

		a.QueueUpdateDraw(func() {
			a.root.chainInfo.SetHeight(height)
			if price != nil {
				a.root.chainInfo.SetEthPrice(*price)
			}
			if gasPrice != nil {
				a.root.chainInfo.SetGasPrice(gasPrice)
			}
		})
	}()

	// load newest transactions
	a.root.home.transactionList.LoadAsync(func() (common.Transactions, error) {
		netType := a.service.GetNetwork().NetType()
		if netType == serv.TypeDevnet {
			return a.service.GetLatestTransactions(100)
		} else {
			return a.service.GetLatestTransactions(1)
		}
	})

	// start syncer
	if err := a.syncer.Start(); err != nil {
		log.Error("Failed to start syncer", "error", err)
		return err
	}

	return nil
}
