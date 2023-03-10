package provider

import (
	"context"
	"encoding/json"
	"math/big"
	"time"

	"github.com/dyng/ramen/internal/common"
	"github.com/dyng/ramen/internal/common/conv"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
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
	Timestamp        uint64          `json:"timeStamp,omitempty"`
}

func (tx *rpcTransaction) UnmarshalJSON(msg []byte) error {
	if err := json.Unmarshal(msg, &tx.tx); err != nil {
		return errors.WithStack(err)
	}
	return json.Unmarshal(msg, &tx.txExtraInfo)
}

func (tx *rpcTransaction) ToTransaction() common.Transaction {
	blockNumer, _ := conv.HexToInt(*tx.BlockNumber)
	return common.WrapTransaction(tx.tx, big.NewInt(blockNumer), tx.From, tx.Timestamp)
}

type rpcBlock struct {
	*types.Header
	rpcBlockBody
}

type rpcBlockBody struct {
	Hash         common.Hash      `json:"hash"`
	Transactions []rpcTransaction `json:"transactions"`
	UncleHashes  []common.Hash    `json:"uncles"`
}

func (b *rpcBlock) UnmarshalJSON(msg []byte) error {
	if err := json.Unmarshal(msg, &b.Header); err != nil {
		return errors.WithStack(err)
	}
	if err := json.Unmarshal(msg, &b.rpcBlockBody); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (b *rpcBlock) ToBlock() *common.Block {
	txns := make(types.Transactions, len(b.Transactions))
	for i, tx := range b.Transactions {
		txns[i] = tx.tx
	}
	return types.NewBlockWithHeader(b.Header).WithBody(txns, []*types.Header{})
}

// forkchainSigner is a signer handles mixed transactions of different chains (often the case of local fork chain)
type forkchainSigner struct {
	chainId *big.Int
	signers map[uint64]types.Signer
}

func newForkchainSigner(chainId *big.Int) types.Signer {
	return &forkchainSigner{
		chainId: chainId,
		signers: make(map[uint64]types.Signer),
	}
}

func (s forkchainSigner) getSigner(tx *types.Transaction) types.Signer {
	signer, ok := s.signers[tx.ChainId().Uint64()]
	if !ok {
		signer = types.NewLondonSigner(tx.ChainId())
		s.signers[tx.ChainId().Uint64()] = signer
	}
	return signer
}

// ChainID implements types.Signer
func (s forkchainSigner) ChainID() *big.Int {
	return s.chainId
}

// Equal implements types.Signer
func (s forkchainSigner) Equal(s2 types.Signer) bool {
	x, ok := s2.(forkchainSigner)
	return ok && x.chainId.Cmp(s.chainId) == 0
}

// Hash implements types.Signer
func (s forkchainSigner) Hash(tx *types.Transaction) common.Hash {
	return s.getSigner(tx).Hash(tx)
}

// Sender implements types.Signer
func (s forkchainSigner) Sender(tx *types.Transaction) (common.Address, error) {
	return s.getSigner(tx).Sender(tx)
}

// SignatureValues implements types.Signer
func (s forkchainSigner) SignatureValues(tx *types.Transaction, sig []byte) (R *big.Int, S *big.Int, V *big.Int, err error) {
	return s.getSigner(tx).SignatureValues(tx, sig)
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
		log.Error("Cannot connect to rpc server", "url", url, "error", errors.WithStack(err))
		common.Exit("Cannot connect to rpc server %s: %v", url, err)
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
			return nil, errors.WithStack(err)
		}
		p.chainId = chainId
		p.signer = newForkchainSigner(chainId)
	}
	return p.chainId, nil
}

func (p *Provider) GetGasPrice() (common.BigInt, error) {
	ctx, cancel := p.createContext()
	defer cancel()
	gasPrice, err := p.client.SuggestGasPrice(ctx)
	return gasPrice, errors.WithStack(err)
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
	code, err := p.client.CodeAt(ctx, addr, nil)
	return code, errors.WithStack(err)
}

func (p *Provider) GetBalance(addr common.Address) (common.BigInt, error) {
	ctx, cancel := p.createContext()
	defer cancel()
	balance, err := p.client.BalanceAt(ctx, addr, nil)
	return balance, errors.WithStack(err)
}

func (p *Provider) GetBlockHeight() (uint64, error) {
	ctx, cancel := p.createContext()
	defer cancel()
	height, err := p.client.BlockNumber(ctx)
	return height, errors.WithStack(err)
}

func (p *Provider) GetBlockByHash(hash common.Hash) (*common.Block, error) {
	ctx, cancel := p.createContext()
	defer cancel()
	block, err := p.client.BlockByHash(ctx, hash)
	return block, errors.WithStack(err)
}

func (p *Provider) GetBlockByNumber(number common.BigInt) (*common.Block, error) {
	ctx, cancel := p.createContext()
	defer cancel()
	block, err := p.client.BlockByNumber(ctx, number)
	return block, errors.WithStack(err)
}

func (p *Provider) BatchBlockByNumber(numberList []common.BigInt) ([]*common.Block, error) {
	size := len(numberList)
	rpcRes := make([]rpcBlock, size)
	reqs := make([]rpc.BatchElem, size)
	for i := range reqs {
		reqs[i] = rpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []any{toBlockNumArg(numberList[i]), true},
			Result: &rpcRes[i],
		}
	}

	ctx, cancel := p.createContext()
	defer cancel()

	err := p.rpcClient.BatchCallContext(ctx, reqs)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := make([]*common.Block, size)
	for i := range result {
		result[i] = rpcRes[i].ToBlock()
	}

	// FIXME: individual request error handling
	return result, nil
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
		return nil, errors.WithStack(err)
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
		return 0, errors.WithStack(err)
	}

	return gasLimit, nil
}

func (p *Provider) CallContract(address common.Address, abi *abi.ABI, method string, args ...any) ([]any, error) {
	// encode calldata
	input, err := abi.Pack(method, args...)
	if err != nil {
		return nil, errors.WithStack(err)
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
		return nil, errors.WithStack(err)
	}

	vals, err := abi.Unpack(method, data)
	if err != nil {
		return nil, errors.WithStack(err)
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
		return common.Hash{}, errors.WithStack(err)
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
		return common.Hash{}, errors.WithStack(err)
	}

	err = p.client.SendTransaction(ctx, signedTx)
	return signedTx.Hash(), errors.WithStack(err)
}

func (p *Provider) SubscribeNewHead(ch chan<- *common.Header) (ethereum.Subscription, error) {
	ctx, cancel := p.createContext()
	defer cancel()
	sub, err := p.client.SubscribeNewHead(ctx, ch)
	return sub, errors.WithStack(err)
}

func (p *Provider) createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), DefaultTimeout)
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	return hexutil.EncodeBig(number)
}
