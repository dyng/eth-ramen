package etherscan

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

const (
	testEtherscanEndpoint = "https://api.etherscan.io/api"
	testEtherscanApiKey = "IQVJUFHSK9SG8SVDK3MKPIJHQR3137GCPQ"

	usdtContractAddress = "0xdAC17F958D2ee523a2206206994597C13D831ec7"
)

func TestAccountTxList_NoError(t *testing.T) {
	// prepare
	ec := NewEtherscanClient(testEtherscanEndpoint, testEtherscanApiKey)

	// process
	txns, err := ec.AccountTxList(common.HexToAddress("0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae"))

	// verify
	assert.NoError(t, err)
	assert.NotEmpty(t, txns, "should not empty")

	txn := txns[0]
	assert.NotNil(t, txn.BlockNumber(), "block number should not be nil")
	assert.NotNil(t, txn.Hash(), "hash should not be nil")
	assert.NotNil(t, txn.From(), "sender should not be nil")
	assert.NotNil(t, txn.To(), "receiver should not be nil")
	assert.NotEmpty(t, txn.Timestamp(), "timestamp should not be nil")
}

func TestGetSourceCode_NoError(t *testing.T) {
	// prepare
	ec := NewEtherscanClient(testEtherscanEndpoint, testEtherscanApiKey)

	// process
	source, abi, err := ec.GetSourceCode(common.HexToAddress(usdtContractAddress))

	// verify
	assert.NoError(t, err)
	assert.NotEmpty(t, source, "source code should not be empty")
	assert.NotNil(t, abi, "abi should not be null")
	assert.Contains(t, abi.Methods, "balanceOf", "should contains method balanceOf")
}
