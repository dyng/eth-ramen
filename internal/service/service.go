package service

import (
	"embed"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/dyng/ramen/internal/common"
	"github.com/dyng/ramen/internal/common/conv"
	conf "github.com/dyng/ramen/internal/config"
	"github.com/dyng/ramen/internal/provider"
	"github.com/dyng/ramen/internal/provider/etherscan"
	gcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/shopspring/decimal"
)

//go:embed data/chains.json
var chainFile embed.FS

var chainMap map[string]Network

func init() {
	bytes, err := chainFile.ReadFile("data/chains.json")
	if err != nil {
		log.Crit("Cannot read chains.json", "error", err)
	}

	var networks []Network
	err = json.Unmarshal(bytes, &networks)
	if err != nil {
		log.Crit("Cannot parse chains.json", "error", err)
	}

	cache := make(map[string]Network)
	for _, n := range networks {
		cache[n.ChainId.String()] = n
	}

	chainMap = cache
}

type Service struct {
	config   *conf.Config
	esclient *etherscan.EtherscanClient
	provider *provider.Provider
}

func NewService(config *conf.Config) *Service {
	service := Service{
		config:   config,
		esclient: etherscan.NewEtherscanClient(config.EtherscanEndpoint(), config.EtherscanApiKey),
		provider: provider.NewProvider(config.Endpoint(), config.Provider),
	}

	return &service
}

// GetProvider returns underlying provider instance.
// Usually you don't need to tackle with provider directly.
func (s *Service) GetProvider() *provider.Provider {
	return s.provider
}

// GetNetwork returns the network that provider is connected to.
func (s *Service) GetNetwork() Network {
	chainId, _ := s.provider.GetNetwork()
	network, ok := chainMap[chainId.String()]
	if !ok {
		return Network{
			Name:    "Unknown",
			Title:   "Unknown",
			ChainId: chainId,
		}
	} else {
		return network
	}
}

// GetBlockHeight returns the current block height.
func (s *Service) GetBlockHeight() (uint64, error) {
	return s.provider.GetBlockHeight()
}

// GetGasPrice returns average gas price of last block.
func (s *Service) GetGasPrice() (common.BigInt, error) {
	return s.provider.GetGasPrice()
}

// GetEthPrice returns ETH price in USD.
func (s *Service) GetEthPrice() (*decimal.Decimal, error) {
	return s.esclient.EthPrice()
}

// GetAccount returns an account of given address.
func (s *Service) GetAccount(address string) (*Account, error) {
	addr := gcommon.HexToAddress(address)

	code, err := s.provider.GetCode(addr)
	if err != nil {
		return nil, err
	}

	return &Account{
		service: s,
		address: addr,
		code:    code,
	}, nil
}

// GetLatestTransactions returns transactions of last n blocks. 
func (s *Service) GetLatestTransactions(n int) (common.Transactions, error) {
	max, err := s.GetBlockHeight()
	if err != nil {
		return nil, err
	}

	min := uint64(1)
	cnt := uint64(n)
	if max > cnt-1 {
		min = max - cnt + 1
	}

	transactions := make([]common.Transaction, 0)
	for i := max; i >= min; i-- {
		block, err := s.provider.GetBlockByNumber(new(big.Int).SetUint64(i))
		if err != nil {
			return transactions, err
		}

		txns, err := s.GetTransactionsByBlock(block)
		if err != nil {
			return transactions, err
		}

		transactions = append(transactions, txns...)
	}

	return transactions, nil
}

// GetTransactionsByBlock returns transactions of given block hash.
func (s *Service) GetTransactionsByBlock(block *common.Block) (common.Transactions, error) {
	signer, err := s.provider.GetSigner()
	if err != nil {
		return nil, err
	}

	txns := make(common.Transactions, block.Transactions().Len())
	for i, tx := range block.Transactions() {
		sender, err := signer.Sender(tx)
		if err != nil {
			return nil, err
		}
		txns[i] = common.WrapTransactionWithBlock(tx, block, &sender)
	}

	return txns, nil
}

// GetTransactionHistory returns transactions related to specified account.
// This method relies on Etherscan API at chains other than local chain.
func (s *Service) GetTransactionHistory(address common.Address) (common.Transactions, error) {
	netType := s.GetNetwork().NetType()
	switch netType {
	case TypeDevnet:
		return s.transactionsByTraverse(address)
	default:
		return s.transactionsByEtherscan(address)
	}
}

func (s *Service) transactionsByTraverse(address common.Address) (common.Transactions, error) {
	candidates, err := s.GetLatestTransactions(100)
	if err != nil {
		return nil, err
	}

	txns := make([]common.Transaction, 0)
	for _, t := range candidates {
		if t.From().String() == address.String() {
			txns = append(txns, t)
		}
		if t.To() != nil && t.To().String() == address.String() {
			txns = append(txns, t)
		}
	}

	return txns, nil
}

func (s *Service) transactionsByEtherscan(address common.Address) (common.Transactions, error) {
	return s.esclient.AccountTxList(address)
}

func (s *Service) transactionsByAlchemy(address common.Address) (common.Transactions, error) {
	hashList := make([]common.Hash, 0)

	// incoming transactions
	params := provider.GetAssetTransfersParams{
		FromAddress: address.Hex(),
		Category:    []string{"external"},
		Order:       "desc",
		MaxCount:    "0x14", // decimal value: 20
	}
	result, err := s.provider.GetAssetTransfers(params)
	if err != nil {
		return nil, err
	}
	transfers := result.Transfers
	for _, tr := range transfers {
		hashList = append(hashList, gcommon.HexToHash(tr.Hash))
	}

	// outgoing transactions
	params = provider.GetAssetTransfersParams{
		ToAddress: address.Hex(),
		Category:  []string{"external"},
		Order:     "desc",
		MaxCount:  "0x14", // decimal value: 20
	}
	result, err = s.provider.GetAssetTransfers(params)
	if err != nil {
		return nil, err
	}
	transfers = result.Transfers
	for _, tr := range transfers {
		hashList = append(hashList, gcommon.HexToHash(tr.Hash))
	}

	txns, err := s.provider.BatchTransactionByHash(hashList)
	if err != nil {
		return nil, err
	}

	return txns, nil
}

// GetContract returns a contract object of given address.
func (s *Service) GetContract(address common.Address) (*Contract, error) {
	account, err := s.GetAccount(address.Hex())
	if err != nil {
		return nil, err
	}

	if account.GetType() != TypeContract {
		return nil, fmt.Errorf("Address %s is not a contract account", address.Hex())
	}

	return s.ToContract(account)
}

// GetSigner returns a signer which can sign transactions
func (s *Service) GetSigner(privateKey string) (*Signer, error) {
	privKey, err := crypto.HexToECDSA(conv.Trim0xPrefix(privateKey))
	if err != nil {
		return nil, err
	}

	// only EOA can have private key
	addr := crypto.PubkeyToAddress(privKey.PublicKey)
	account := &Account{
		service: s,
		address: addr,
	}

	signer := &Signer{
		Account:    account,
		PrivateKey: privKey,
	}
	return signer, nil
}

// ToContract upgrade an account object to a contract.
func (s *Service) ToContract(account *Account) (*Contract, error) {
	if account.GetType() != TypeContract {
		return nil, fmt.Errorf("Address %s is not a contract account", account.address.Hex())
	}

	if s.GetNetwork().NetType() == TypeDevnet {
		// in case of devnet, return a contract skeleton
		return &Contract{
			Account: account,
		}, nil
	}

	source, abi, err := s.esclient.GetSourceCode(account.address)
	if err != nil {
		return nil, err
	}

	if source == "" {
		return &Contract{
			Account: account,
		}, nil
	}

	return &Contract{
		Account: account,
		abi:     abi,
		source:  source,
	}, nil
}
