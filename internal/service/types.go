package service

import (
	"math/big"
	"strings"
)

const (
	// TypeMainnet is the Ethereum Mainnet
	TypeMainnet = "mainnet"
	// TypeTestnet is all the testnets (Ropsten, Rinkeby, Goerli etc.)
	TypeTestnet = "testnet"
	// TypeDevnet is a local network for development purpose (Hardhat, Ganeche etc.)
	TypeDevnet = "devnet"
	// TypeUnknown is a unknown network
	TypeUnknown = "unknown"
)

type Network struct {
	Name    string `json:"name"`
	Title   string `json:"title"`
	ChainId *big.Int `json:"chainId"`
}

func (n Network) NetType() string {
	if n.Name == "Ethereum Mainnet" {
		return TypeMainnet
	}

	if strings.Contains(n.Title, "Testnet") {
		return TypeTestnet
	}

	if n.ChainId.String() == "1337" || n.ChainId.String() == "31337" {
		return TypeDevnet
	}

	return TypeUnknown
}

