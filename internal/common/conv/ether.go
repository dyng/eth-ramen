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
func FromEther(n *big.Float) *big.Int {
	i, _ := new(big.Float).Mul(n, big.NewFloat(params.Ether)).Int(nil)
	return i
}

// FromGwei converts values in Gwei to wei.
func FromGwei(n *big.Float) *big.Int {
	i, _ := new(big.Float).Mul(n, big.NewFloat(params.GWei)).Int(nil)
	return i
}
