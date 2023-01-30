package common

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Transaction represents an Ethereum transaction.
type Transaction interface {
	BlockNumber() BigInt

	Hash() Hash

	From() *Address

	To() *Address

	Value() BigInt

	Timestamp() uint64
}

// TxnRequest represents a transaction to be submitted for execution
type TxnRequest struct {
	PrivateKey *ecdsa.PrivateKey
	To         *Address
	Value      BigInt
	Data       []byte
	GasLimit   uint64
	GasPrice   BigInt
}

// WrappedTransaction is a wrapper around geth Transaction for convenience
type WrappedTransaction struct {
	*types.Transaction
	from        *common.Address
	blockNumber BigInt
	timestamp   uint64
}

type Transactions = []Transaction

func WrapTransaction(txn *types.Transaction, blockNumber BigInt, from *common.Address, timestamp uint64) Transaction {
	return &WrappedTransaction{
		Transaction: txn,
		from:        from,
		blockNumber: blockNumber,
		timestamp:   timestamp,
	}
}

func WrapTransactionWithBlock(txn *types.Transaction, block *types.Block, sender *common.Address) Transaction {
	return &WrappedTransaction{
		Transaction: txn,
		from:        sender,
		blockNumber: block.Number(),
		timestamp:   block.Time(),
	}
}

func (t *WrappedTransaction) From() *Address {
	return t.from
}

func (t *WrappedTransaction) BlockNumber() BigInt {
	return t.blockNumber
}

func (t *WrappedTransaction) Timestamp() uint64 {
	return t.timestamp
}
