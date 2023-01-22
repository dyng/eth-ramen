package common

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// BigInt is used everywhere in Ethereum
type BigInt = *big.Int

// Address is an alias for geth Address
type Address = common.Address

// Block is an alias for geth Block
type Block = types.Block

// Hash is an alias for geth Hash
type Hash = common.Hash

// Header is an alias for geth Header
type Header = types.Header
