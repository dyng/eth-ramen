package provider

import (
	"math/big"
	"strings"
	"testing"

	"github.com/dyng/ramen/internal/common"
	"github.com/ethereum/go-ethereum/accounts/abi"
	gcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

const (
	testAlchemyEndpoint = "wss://eth-mainnet.g.alchemy.com/v2/1DYmd-KT-4evVd_-O56p5HTgk2t5cuVu"

	usdtContractAddress = "0xdAC17F958D2ee523a2206206994597C13D831ec7"
	usdtFunctionABI     = "[{\"constant\":true,\"inputs\":[{\"name\":\"who\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"
)

func TestBatchTransactionByHash_NoError(t *testing.T) {
	// prepare
	provider := NewProvider(testAlchemyEndpoint, ProviderAlchemy)

	// process
	hashList := []common.Hash{
		gcommon.HexToHash("0xc2a5c78171f96e1268035ee8c90436dc6945a73b03a4970a6c38f1635a6a1bd2"),
		gcommon.HexToHash("0x5e9a9c54899325819a04e591d1930bee53b8798de3d26f438e9beba89de2fafa"),
		gcommon.HexToHash("0xf53efe987616ecf1b1178de1dc8ae58946f119003774abb132c3cd9c3fccb762"),
	}
	txns, err := provider.BatchTransactionByHash(hashList)

	// verify
	assert.NoError(t, err)
	assert.Len(t, txns, 3)
	for _, txn := range txns {
		assert.NotNil(t, txn.BlockNumber(), "block number should not be nil")
		assert.NotNil(t, txn.Hash(), "hash should not be nil")
		assert.NotNil(t, txn.From(), "sender should not be nil")
		assert.NotNil(t, txn.To(), "receiver should not be nil")
	}
}

func TestBatchBlockByNumber_NoError(t *testing.T) {
	// prepare
	provider := NewProvider(testAlchemyEndpoint, ProviderAlchemy)

	// process
	numberList := []common.BigInt{
		big.NewInt(16748002),
		big.NewInt(16748001),
		big.NewInt(16748000),
	}
	blocks, err := provider.BatchBlockByNumber(numberList)

	// verify
	assert.NoError(t, err)
	assert.Len(t, blocks, 3)
	for _, block := range blocks {
		assert.NotNil(t, block.Number(), "block number should not be nil")
		assert.NotNil(t, block.Hash(), "hash should not be nil")
		assert.NotEmpty(t, block.Transactions(), "transactions should not be empty")
	}
}

func TestCallContract_NoError(t *testing.T) {
	// prepare
	provider := NewProvider(testAlchemyEndpoint, ProviderAlchemy)
	usdtABI, _ := abi.JSON(strings.NewReader(usdtFunctionABI))
	usdtAddr := gcommon.HexToAddress(usdtContractAddress)
	argAddr := gcommon.HexToAddress("0x759B7e31E6411AB92CF382b3d4733D98134052a7")

	// process
	result, err := provider.CallContract(usdtAddr, &usdtABI, "balanceOf", argAddr)

	// verify
	assert.NoError(t, err)
	balance := result[0].(common.BigInt)
	assert.Equal(t, 1, balance.Cmp(big.NewInt(100)), "balance should be greater than 100")
}

func TestGetGasPrice_NoError(t *testing.T) {
	// prepare
	provider := NewProvider(testAlchemyEndpoint, ProviderAlchemy)

	// process
	result, err := provider.GetGasPrice()

	// verify
	assert.NoError(t, err)
	assert.Equal(t, 1, result.Cmp(big.NewInt(100)), "gas price should be greater than 100")
}
