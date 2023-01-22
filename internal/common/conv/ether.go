package conv

import (
	"math/big"

	"github.com/ethereum/go-ethereum/params"
)

func ToEther(wei *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(params.Ether))
}

