package service

import (
	"github.com/dyng/ramen/internal/common"
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

// GetBalance returns current balance of this account.
func (a *Account) GetBalance() (common.BigInt, error) {
	return a.service.provider.GetBalance(a.address)
}

// GetTransactions returns transactions of this account.
func (a *Account) GetTransactions() (common.Transactions, error) {
	return a.service.GetTransactionHistory(a.address)
}

// AsContract upgrade this account object to a contract.
func (a *Account) AsContract() (*Contract, error) {
	return a.service.ToContract(a)
}
