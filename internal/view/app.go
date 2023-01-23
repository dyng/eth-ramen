package view

import (
	"github.com/asaskevich/EventBus"
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
	syncer := serv.NewSyncer(app.service.GetProvider(), app.eventBus)
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
	conf := a.config

	// update network
	a.root.chainInfo.SetNetwork(conf.Network)

	// update block height
	height, err := a.service.GetBlockHeight()
	if err != nil {
		log.Error("Failed to fetch block height", "error", err)
	} else {
		a.root.chainInfo.SetHeight(height)
	}

	// show latest transactions
	a.root.home.transactionList.LoadAsync(a.service.GetLatestTransactions)

	// start syncer
	if err := a.syncer.Start(); err != nil {
		log.Error("Failed to start syncer", "error", err)
		return err
	}

	return nil
}
