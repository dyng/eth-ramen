package config

import (
	"fmt"

	"github.com/dyng/ramen/internal/view/style"
)

const (
	// DefaultProvider represents the Ethereum provider
	DefaultProvider = "local"

	// DefaultNetwork represents the chain we connect to
	DefaultNetwork = "mainnet"
)

type Config struct {
	Provider        string
	Network         string
	ApiKey          string
	EtherscanApiKey string
}

func NewConfig() *Config {
	return &Config{
		Provider:        DefaultProvider,
		Network:         DefaultNetwork,
		// FIXME: delete keys
		ApiKey:          "1DYmd-KT-4evVd_-O56p5HTgk2t5cuVu",
		EtherscanApiKey: "IQVJUFHSK9SG8SVDK3MKPIJHQR3137GCPQ",
	}
}

func (c *Config) Endpoint() string {
	// config api key
	apiKey := c.ApiKey

	// config network
	network := c.Network
	if network == "homestead" {
		network = "mainnet"
	}

	switch c.Provider {
	case "local":
		return "ws://localhost:8545"
	case "alchemy":
		return fmt.Sprintf("wss://eth-%s.alchemyapi.io/v2/%s", network, apiKey)
	default:
		return ""
	}
}

func (c *Config) EtherscanEndpoint() string {
	if c.Network == "mainnet" {
		return fmt.Sprintf("https://api.etherscan.io/api")
	} else {
		return fmt.Sprintf("https://api-%s.etherscan.io/api", c.Network)
	}
}

func (c *Config) Style() *style.Style {
	return style.Ethereum
}
