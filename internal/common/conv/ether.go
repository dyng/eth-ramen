package conv

import (
	"math/big"

	"github.com/ethereum/go-ethereum/params"
)

// ToEther converts values in wei to ether.
func ToEther(wei *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(params.Ether))
}

// ToGwei converts values in wei to Gwei.
func ToGwei(wei *big.Int) *big.Int {
	return new(big.Int).Quo(wei, big.NewInt(params.GWei))
}

// FromEther converts values in Ether to wei.
func FromEther(n int64) *big.Int {
	return new(big.Int).Mul(big.NewInt(n), big.NewInt(params.Ether))
}

// FromGwei converts values in Gwei to wei.
func FromGwei(n int64) *big.Int {
	return new(big.Int).Mul(big.NewInt(n), big.NewInt(params.GWei))
}

