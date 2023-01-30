package service

import (
	"testing"

	"github.com/dyng/ramen/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestGetNetwork(t *testing.T) {
	// prepare
	serv := NewTestService()

	// process
	network := serv.GetNetwork()

	// verify
	assert.Equal(t, "31337", network.ChainId.String(), "chain id should be 31337")
	assert.Equal(t, "GoChain Testnet", network.Name, "chain name should be GoChain Testnet")
}

func TestGetSigner(t *testing.T) {
	// prepare
	serv := NewTestService()

	// process
	privateKey := "0xde9be858da4a475276426320d5e9262ecfc3ba460bfac56360bfa6c4c28b4ee0"
	signer, err := serv.GetSigner(privateKey)

	// verify
	assert.NoError(t, err)
	assert.Equal(t, signer.GetAddress().Hex(), "0xdD2FD4581271e230360230F9337D5c0430Bf44C0", "signer's account address should be correct")
	assert.NotNil(t, signer.PrivateKey, "signer should have private key")
}

func NewTestService() *Service {
	config := &config.Config{
		DebugMode: true,
		Provider:  "local",
		Network:   "mainnet",
	}
	return NewService(config)
}
