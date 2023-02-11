package service

import (
	"math/big"
	"sync/atomic"

	"github.com/dyng/ramen/internal/common"
	"github.com/ethereum/go-ethereum/log"
)

const (
	// TypeWallet is an EOA account.
	TypeWallet AccountType = "Wallet"
	// TypeContract is a SmartContract account.
	TypeContract = "Contract"
)

// AccountType represents two types of Etheruem's account: EOA and SmartContract.
type AccountType string

func (at AccountType) String() string {
	return string(at)
}

// Account represents an account of Etheruem network.
type Account struct {
	service *Service
	address common.Address
	balance atomic.Pointer[big.Int]
	code    []byte // byte code of this account, nil if account is an EOA.
}

// GetAddress returns address of this account.
func (a *Account) GetAddress() common.Address {
	return a.address
}

// GetType returns type of this account, either Wallet or Contract.
func (a *Account) GetType() AccountType {
	if len(a.code) == 0 {
		return TypeWallet
	} else {
		return TypeContract
	}
}

// IsContract returns true if this account is a smart contract.
func (a *Account) IsContract() bool {
	return a.GetType() == TypeContract
}

// AsContract upgrade this account object to a contract.
func (a *Account) AsContract() (*Contract, error) {
	return a.service.ToContract(a)
}

// GetBalance returns cached balance of this account.
func (a *Account) GetBalance() common.BigInt {
	// FIXME: race condition
	if a.balance.Load() == nil {
		bal, err := a.GetBalanceForce()
		if err != nil {
			log.Error("Failed to fetch balance", "address", a.address, "error", err)
		}
		return bal
	} else {
		return a.balance.Load()
	}
}

// GetBalanceForce will query for current account's balance, store it in cache and return.
func (a *Account) GetBalanceForce() (common.BigInt, error) {
	bal, err := a.service.provider.GetBalance(a.address)
	if err == nil {
		a.balance.Swap(bal)
	} else {
		bal = big.NewInt(0) // use 0 as fallback value
	}
	return bal, err
}

// UpdateBalance will update cache of current account's balance
func (a *Account) UpdateBalance() bool {
	_, err := a.GetBalanceForce()
	return err == nil
}

// ClearCache will clear cached balance
func (a *Account) ClearCache() {
	a.balance.Store(nil)
}

// GetTransactions returns transactions of this account.
func (a *Account) GetTransactions() (common.Transactions, error) {
	return a.service.GetTransactionHistory(a.address)
}
