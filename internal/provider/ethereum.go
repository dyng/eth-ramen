package provider

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"time"

	"github.com/dyng/ramen/internal/common"
	"github.com/dyng/ramen/internal/common/conv"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
)

var (
	ErrProviderNotSupport = errors.New("provider does not support this vendor-specific api")
)

const (
	// ProviderLocal represents a local node provider such as Geth, Hardhat etc.
	ProviderLocal string = "local"
	// ProviderAlchemy represents blockchain provider ProviderAlchemy (https://www.alchemy.com/)
	ProviderAlchemy = "alchemy"

	// DefaultTimeout is the default value for request timeout
	DefaultTimeout = 30 * time.Second
)

type rpcTransaction struct {
	tx *types.Transaction
	txExtraInfo
}

type txExtraInfo struct {
	BlockNumber      *string         `json:"blockNumber,omitempty"`
	BlockHash        *common.Hash    `json:"blockHash,omitempty"`
	From             *common.Address `json:"from,omitempty"`
	TransactionIndex uint            `json:"transactionIndex,omitempty"`
	Timestamp        uint64          `json:"timeStamp,omitempty"`
}

func (tx *rpcTransaction) UnmarshalJSON(msg []byte) error {
	if err := json.Unmarshal(msg, &tx.tx); err != nil {
		return err
	}
	return json.Unmarshal(msg, &tx.txExtraInfo)
}

func (tx *rpcTransaction) ToTransaction() common.Transaction {
	blockNumer, _ := conv.HexToInt(*tx.BlockNumber)
	return common.WrapTransaction(tx.tx, big.NewInt(blockNumer), tx.From, tx.Timestamp)
}

type Provider struct {
	url          string
	providerType string
	client       *ethclient.Client
	rpcClient    *rpc.Client

	// cache
	chainId common.BigInt
	signer  types.Signer
}

// NewProvider returns
func NewProvider(url string, providerType string) *Provider {
	p := &Provider{
		url:          url,
		providerType: providerType,
	}

	rpcClient, err := rpc.Dial(url)
	if err != nil {
		log.Crit("Cannot connect to rpc server", "url", url, "error", err)
	}

	p.rpcClient = rpcClient
	p.client = ethclient.NewClient(rpcClient)

	return p
}

func (p *Provider) GetType() string {
	return p.providerType
}

func (p *Provider) GetNetwork() (common.BigInt, error) {
	ctx, cancel := p.createContext()
	defer cancel()

	if p.chainId == nil {
		chainId, err := p.client.NetworkID(ctx)
		if err != nil {
			return nil, err
		}
		p.chainId = chainId
		p.signer = types.NewLondonSigner(chainId)
	}
	return p.chainId, nil
}

func (p *Provider) GetGasPrice() (common.BigInt, error) {
	ctx, cancel := p.createContext()
	defer cancel()
	return p.client.SuggestGasPrice(ctx)
}

func (p *Provider) GetSigner() (types.Signer, error) {
	_, err := p.GetNetwork()
	if err != nil {
		return nil, err
	}
	return p.signer, nil
}

func (p *Provider) GetCode(addr common.Address) ([]byte, error) {
	ctx, cancel := p.createContext()
	defer cancel()
	return p.client.CodeAt(ctx, addr, nil)
}

func (p *Provider) GetBalance(addr common.Address) (common.BigInt, error) {
	ctx, cancel := p.createContext()
	defer cancel()
	return p.client.BalanceAt(ctx, addr, nil)
}

func (p *Provider) GetBlockHeight() (uint64, error) {
	ctx, cancel := p.createContext()
	defer cancel()
	return p.client.BlockNumber(ctx)
}

func (p *Provider) GetBlockByHash(hash common.Hash) (*common.Block, error) {
	ctx, cancel := p.createContext()
	defer cancel()
	return p.client.BlockByHash(ctx, hash)
}

func (p *Provider) GetBlockByNumber(number common.BigInt) (*common.Block, error) {
	ctx, cancel := p.createContext()
	defer cancel()
	return p.client.BlockByNumber(ctx, number)
}

func (p *Provider) BatchTransactionByHash(hashList []common.Hash) (common.Transactions, error) {
	size := len(hashList)
	rpcRes := make([]rpcTransaction, size)
	reqs := make([]rpc.BatchElem, size)
	for i := range reqs {
		reqs[i] = rpc.BatchElem{
			Method: "eth_getTransactionByHash",
			Args:   []any{hashList[i]},
			Result: &rpcRes[i],
		}
	}

	ctx, cancel := p.createContext()
	defer cancel()

	err := p.rpcClient.BatchCallContext(ctx, reqs)
	if err != nil {
		return nil, err
	}

	result := make(common.Transactions, size)
	for i := range result {
		result[i] = rpcRes[i].ToTransaction()
	}

	// FIXME: individual request error handling
	return result, nil
}

func (p *Provider) EstimateGas(address common.Address, from common.Address, input []byte) (uint64, error) {
	// build call message
	msg := ethereum.CallMsg{
		From: from,
		To:   &address,
		Data: input,
	}

	ctx, cancel := p.createContext()
	defer cancel()

	gasLimit, err := p.client.EstimateGas(ctx, msg)
	if err != nil {
		return 0, err
	}

	return gasLimit, nil
}

func (p *Provider) CallContract(address common.Address, abi *abi.ABI, method string, args ...any) ([]any, error) {
	// encode calldata
	input, err := abi.Pack(method, args...)
	if err != nil {
		return nil, err
	}

	// build call message
	msg := ethereum.CallMsg{
		To:   &address,
		Data: input,
	}

	ctx, cancel := p.createContext()
	defer cancel()

	data, err := p.client.CallContract(ctx, msg, nil)
	if err != nil {
		return nil, err
	}

	vals, err := abi.Unpack(method, data)
	if err != nil {
		return nil, err
	}

	return vals, nil
}

func (p *Provider) SendTransaction(txnReq *common.TxnRequest) (common.Hash, error) {
	ctx, cancel := p.createContext()
	defer cancel()

	key := txnReq.PrivateKey
	from := crypto.PubkeyToAddress(key.PublicKey)

	// fetch the next nonce
	nonce, err := p.client.PendingNonceAt(ctx, from)
	if err != nil {
		return common.Hash{}, err
	}

	txn := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: txnReq.GasPrice,
		Gas:      txnReq.GasLimit,
		To:       txnReq.To,
		Value:    txnReq.Value,
		Data:     txnReq.Data,
	})

	signer, err := p.GetSigner()
	if err != nil {
		return common.Hash{}, err
	}

	signedTx, err := types.SignTx(txn, signer, key)
	if err != nil {
		return common.Hash{}, err
	}

	err = p.client.SendTransaction(ctx, signedTx)
	return signedTx.Hash(), err
}

func (p *Provider) SubscribeNewHead(ch chan<- *common.Header) (ethereum.Subscription, error) {
	ctx, cancel := p.createContext()
	defer cancel()
	return p.client.SubscribeNewHead(ctx, ch)
}

func (p *Provider) createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), DefaultTimeout)
}
