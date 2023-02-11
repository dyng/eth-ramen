package service

import (
	"time"

	"github.com/dyng/ramen/internal/common"
)

func (s *Service) SetCache(address common.Address, accountType AccountType, value any, expiration time.Duration) {
	s.cache.Set(s.cacheKey(address, accountType), value, expiration)
}

func (s *Service) GetCache(address common.Address, accountType AccountType) (any, bool) {
	return s.cache.Get(s.cacheKey(address, accountType))
}

func (s *Service) cacheKey(address common.Address, accountType AccountType) string {
	chainId := s.GetNetwork().ChainId
	return chainId.String() + ":" + address.Hex() + ":" + accountType.String()
}

