package etherscan

import (
	"encoding/json"
	"math/big"

	"github.com/dyng/ramen/internal/common"
	"github.com/dyng/ramen/internal/common/conv"
	"github.com/pkg/errors"
)

type resMessage struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
}

type esTransaction struct {
	blockNumber      common.BigInt
	timeStamp        uint64
	hash             common.Hash
	nonce            uint64
	blockHash        common.Hash
	transactionIndex uint
	from             *common.Address
	to               *common.Address
	value            common.BigInt
	gas              uint64
	gasPrice         common.BigInt
	data             []byte
}

// BlockNumber implements common.Transaction
func (t *esTransaction) BlockNumber() common.BigInt {
	return t.blockNumber
}

// From implements common.Transaction
func (t *esTransaction) From() *common.Address {
	return t.from
}

// Hash implements common.Transaction
func (t *esTransaction) Hash() common.Hash {
	return t.hash
}

// Timestamp implements common.Transaction
func (t *esTransaction) Timestamp() uint64 {
	return t.timeStamp
}

// To implements common.Transaction
func (t *esTransaction) To() *common.Address {
	return t.to
}

// Value implements common.Transaction
func (t *esTransaction) Value() common.BigInt {
	return t.value
}

// Data implements common.Transaction
func (t *esTransaction) Data() []byte {
	return t.data
}

type txJSON struct {
	BlockNumber      int64           `json:"blockNumber,string"`
	TimeStamp        uint64          `json:"timeStamp,string"`
	Hash             common.Hash     `json:"hash"`
	Nonce            uint64          `json:"nonce,string"`
	BlockHash        common.Hash     `json:"blockHash"`
	TransactionIndex uint            `json:"transactionIndex,string"`
	From             *common.Address `json:"from"`
	To               *common.Address `json:"to"`
	Value            string          `json:"value"`
	Gas              uint64          `json:"gas,string"`
	GasPrice         string          `json:"gasPrice"`
	IsError          int64           `json:"isError,string"`
	TxReceiptStatus  int64           `json:"txreceipt_status,string"`
	Input            string          `json:"input"`
}

func (t *esTransaction) UnmarshalJSON(input []byte) error {
	var tx txJSON
	err := json.Unmarshal(input, &tx)
	if err != nil {
		return errors.WithStack(err)
	}

	t.blockNumber = big.NewInt(tx.BlockNumber)
	t.timeStamp = tx.TimeStamp
	t.hash = tx.Hash
	t.nonce = tx.Nonce
	t.blockHash = tx.BlockHash
	t.transactionIndex = tx.TransactionIndex
	t.from = tx.From
	t.to = tx.To
	t.gas = tx.Gas

	bi, ok := new(big.Int).SetString(tx.Value, 10)
	if !ok {
		return errors.Errorf("cannot convert value %s to big.Int", tx.Value)
	}
	t.value = bi

	bi, ok = new(big.Int).SetString(tx.GasPrice, 10)
	if !ok {
		return errors.Errorf("cannot convert value %s to big.Int", tx.GasPrice)
	}
	t.gasPrice = bi

	bs, err := conv.HexToBytes(tx.Input)
	if err != nil {
		return errors.WithStack(err)
	}
	t.data = bs

	return nil
}

type contractJSON struct {
	SourceCode   string `json:"SourceCode"`
	ABI          string `json:"ABI"`
	ContractName string `json:"ContractName"`
}

type ethpriceJSON struct {
	EthBtc          string `json:"ethbtc"`
	EthBtcTimestamp string `json:"ethbtc_timestamp"`
	EthUsd          string `json:"ethusd"`
	EthUsdTimestamp string `json:"ethusd_timestamp"`
}
