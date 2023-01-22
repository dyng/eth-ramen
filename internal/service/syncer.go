package service

import (
	"errors"
	"sync"

	"github.com/asaskevich/EventBus"
	"github.com/dyng/ramen/internal/common"
	"github.com/dyng/ramen/internal/provider"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/log"
)

type Syncer struct {
	*sync.Mutex

	started  bool
	provider *provider.Provider
	eventBus EventBus.Bus
	chBlock  chan *common.Header
	ethSub   ethereum.Subscription
}

func NewSyncer(provider *provider.Provider, eventBus EventBus.Bus) *Syncer {
	return &Syncer{
		Mutex:    &sync.Mutex{},
		started:  false,
		provider: provider,
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

	sub, err := s.provider.SubscribeNewHead(s.chBlock)
	if err != nil {
		return err
	}
	s.ethSub = sub

	// start syncing
	go s.sync()

	return nil
}

func (s *Syncer) sync() {
	for {
		select {
		// new block event
		case newHeader := <-s.chBlock:
			log.Debug("Received new block header", "hash", newHeader.Hash(),
				"number", newHeader.Number)

			block, err := s.provider.GetBlockByHash(newHeader.Hash())
			if err != nil {
				continue
			}

			s.eventBus.Publish(TopicNewBlock, block)
		}
	}
}
