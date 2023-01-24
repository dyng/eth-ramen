package service

import (
	"testing"

	"github.com/dyng/ramen/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestGetNetwork(t *testing.T) {
	// prepare
	conf := ConfigForTesting()
	serv := NewService(conf)

	// process
	network := serv.GetNetwork()

	// verify
	assert.Equal(t, "1", network.ChainId, "chain id should be 1")
	assert.Equal(t, "Mainnet", network.Name, "chain name should be mainnet")
}

func ConfigForTesting() *config.Config {
	return config.NewConfig()
}
