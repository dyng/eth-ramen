package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/dyng/ramen/internal/view/style"
)

var (
	DefaultProvider   = "alchemy"
	DefaultNetwork    = "mainnet"
	DefaultConfigFile = os.Getenv("HOME") + "/.ramen.json"
)

type configJSON struct {
	Provider        *string `json:"provider,omitempty"`
	ApiKey          *string `json:"apikey,omitempty"`
	EtherscanApiKey *string `json:"etherscanApikey,omitempty"`
}

type Config struct {
	DebugMode       bool
	ConfigFile      string
	Provider        string
	Network         string
	ApiKey          string
	EtherscanApiKey string
}

func NewConfig() *Config {
	return &Config{}
}

// ParseConfig extract config file location from Config struct, read and parse
// it, then overwrite Config struct in place.
func ParseConfig(config *Config) error {
	// if config file does not exist, ignore
	path := config.ConfigFile
	bytes, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}

	// read and parse config file
	configJson := new(configJSON)
	err = json.Unmarshal(bytes, &configJson)
	if err != nil {
		return err
	}

	// overwrite configurations
	if configJson.Provider != nil {
		config.Provider = *configJson.Provider
	}
	if configJson.ApiKey != nil {
		config.ApiKey = *configJson.ApiKey
	}
	if configJson.EtherscanApiKey != nil {
		config.EtherscanApiKey = *configJson.EtherscanApiKey
	}

	return nil
}

func (c *Config) Endpoint() string {
	// config api key
	apiKey := c.ApiKey

	// config network
	switch c.Provider {
	case "local":
		return "ws://localhost:8545"
	case "alchemy":
		return fmt.Sprintf("wss://eth-%s.alchemyapi.io/v2/%s", strings.ToLower(c.Network), apiKey)
	default:
		return ""
	}
}

func (c *Config) EtherscanEndpoint() string {
	if c.Network == "mainnet" {
		return fmt.Sprintf("https://api.etherscan.io/api")
	} else {
		return fmt.Sprintf("https://api-%s.etherscan.io/api", strings.ToLower(c.Network))
	}
}

func (c *Config) Style() *style.Style {
	return style.Ethereum
}
