package service

import (
	"embed"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/dyng/ramen/internal/common"
	conf "github.com/dyng/ramen/internal/config"
	"github.com/dyng/ramen/internal/provider"
	"github.com/dyng/ramen/internal/provider/etherscan"
	gcommon "github.com/ethereum/go-ethereum/common"
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

func (s *Service) GetProvider() *provider.Provider {
	return s.provider
}

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

func (s *Service) GetBlockHeight() (uint64, error) {
	return s.provider.GetBlockHeight()
}

func (s *Service) GetGasPrice() (common.BigInt, error) {
	return s.provider.GetGasPrice()
}

func (s *Service) GetEthPrice() (*decimal.Decimal, error) {
	return s.esclient.EthPrice()
}

func (s *Service) GetAccount(address string) (*Account, error) {
	addr := gcommon.HexToAddress(address)
	log.Debug("Try to fetch account", "address", address)

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

func (s *Service) GetLatestTransactions(blockCnt int) (common.Transactions, error) {
	max, err := s.GetBlockHeight()
	if err != nil {
		return nil, err
	}

	min := uint64(1)
	cnt := uint64(blockCnt)
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

func (s *Service) GetTransactionHistory(addr common.Address) (common.Transactions, error) {
	netType := s.GetNetwork().NetType()
	switch netType {
	case TypeDevnet:
		return s.transactionsByTraverse(addr)
	default:
		return s.transactionsByEtherscan(addr)
	}
}

func (s *Service) transactionsByTraverse(addr common.Address) (common.Transactions, error) {
	candidates, err := s.GetLatestTransactions(100)
	if err != nil {
		return nil, err
	}

	txns := make([]common.Transaction, 0)
	for _, t := range candidates {
		if t.From().String() == addr.String() || t.To().String() == addr.String() {
			txns = append(txns, t)
		}
	}

	return txns, nil
}

func (s *Service) transactionsByEtherscan(addr common.Address) (common.Transactions, error) {
	return s.esclient.AccountTxList(addr)
}

func (s *Service) transactionsByAlchemy(addr common.Address) (common.Transactions, error) {
	hashList := make([]common.Hash, 0)

	// incoming transactions
	params := provider.GetAssetTransfersParams{
		FromAddress: addr.Hex(),
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
		ToAddress: addr.Hex(),
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

func (s *Service) GetContract(addr common.Address) (*Contract, error) {
	account, err := s.GetAccount(addr.Hex())
	if err != nil {
		return nil, err
	}

	if account.GetType() != TypeContract {
		return nil, fmt.Errorf("Address %s is not a contract account", addr.Hex())
	}

	return s.ToContract(account)
}

func (s *Service) ToContract(account *Account) (*Contract, error) {
	if account.GetType() != TypeContract {
		return nil, fmt.Errorf("Address %s is not a contract account", account.address.Hex())
	}

	// FIXME: support local chain
	source, abi, err := s.esclient.GetSourceCode(account.address)
	if err != nil {
		return nil, err
	}

	return &Contract{
		Account: account,
		abi:     abi,
		source:  source,
	}, nil
}
