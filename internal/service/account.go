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

type Account struct {
	service *Service
	address common.Address
	code    []byte
}

func (a *Account) GetAddress() common.Address {
	return a.address
}

func (a *Account) GetType() AccountType {
	if len(a.code) == 0 {
		return TypeWallet
	} else {
		return TypeContract
	}
}

func (a *Account) IsContract() bool {
	return a.GetType() == TypeContract
}

func (a *Account) GetBalance() (common.BigInt, error) {
	return a.service.provider.GetBalance(a.address)
}

func (a *Account) GetTransactions() (common.Transactions, error) {
	return a.service.GetTransactionHistory(a.address)
}

func (a *Account) AsContract() (*Contract, error) {
	return a.service.ToContract(a)
}
