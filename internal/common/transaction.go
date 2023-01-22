package common

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Transaction interface {
	BlockNumber() BigInt

	Hash() Hash

	From() *Address

	To() *Address

	Value() BigInt

	Timestamp() uint64
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
