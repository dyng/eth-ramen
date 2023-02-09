package service

import (
	"time"

	"github.com/dyng/ramen/internal/common"
)

func (s *Service) cacheKey(address common.Address) string {
	chainId := s.GetNetwork().ChainId
	return chainId.String() + ":" + address.Hex()
}

func (s *Service) PutCache(address common.Address, value any, expiration time.Duration) {
	s.cache.Set(s.cacheKey(address), value, expiration)
}

func (s *Service) GetCache(address common.Address) (any, bool) {
	return s.cache.Get(s.cacheKey(address))
}
