package service

import (
	"errors"
	"math/big"
	"sync"
	"time"

	"github.com/asaskevich/EventBus"
	"github.com/dyng/ramen/internal/common"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/log"
	"github.com/shopspring/decimal"
)

const (
	// TopicNewBlock is the topic about received new blocks
	TopicNewBlock = "service:newBlock"
	// TopicNewChainData is the topic about latest chain data (ether price, gas price, etc.)
	TopicNewChainData = "service:newChainData"

	// UpdatePeriod is the time duration between two updates
	UpdatePeriod = 5 * time.Second
)

type ChainData struct {
	Price    *decimal.Decimal
	GasPrice *big.Int
}

type Syncer struct {
	*sync.Mutex

	started  bool
	service  *Service
	eventBus EventBus.Bus
	ticker   *time.Ticker
	chBlock  chan *common.Header
	ethSub   ethereum.Subscription
}

func NewSyncer(service *Service, eventBus EventBus.Bus) *Syncer {
	return &Syncer{
		Mutex:    &sync.Mutex{},
		started:  false,
		service:  service,
		eventBus: eventBus,
		chBlock:  make(chan *common.Header),
	}
}

func (s *Syncer) Start() error {
	s.Lock()
	defer s.Unlock()

	if s.started {
		return errors.New("syncer is already started")
	}
	s.started = true

	// subscribe to new blocks
	sub, err := s.service.GetProvider().SubscribeNewHead(s.chBlock)
	if err != nil {
		return err
	}
	s.ethSub = sub

	// start ticker for periodic update
	s.ticker = time.NewTicker(UpdatePeriod)

	// start syncing
	go s.sync()

	return nil
}

func (s *Syncer) sync() {
	for {
		select {
		case newHeader := <-s.chBlock:
			log.Debug("Received new block header", "hash", newHeader.Hash(),
				"number", newHeader.Number)

			block, err := s.service.GetProvider().GetBlockByHash(newHeader.Hash())
			if err != nil {
				log.Error("Failed to fetch block by hash", "hash", newHeader.Hash(), "error", err)
				continue
			}

			s.eventBus.Publish(TopicNewBlock, block)
		case tick := <-s.ticker.C:
			log.Debug("Process periodic synchronization", "tick", tick)

			price, err := s.service.GetEthPrice()
			if err != nil {
				log.Error("Failed to fetch ether's price", "error", err)
			}

			gasPrice, err := s.service.GetGasPrice()
			if err != nil {
				log.Error("Failed to fetch gas price", "error", err)
			}

			data := &ChainData{
				Price:    price,
				GasPrice: gasPrice,
			}
			s.eventBus.Publish(TopicNewChainData, data)
		}
	}
}
