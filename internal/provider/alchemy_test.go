package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAssetTransfers_NoError(t *testing.T) {
	// prepare
	provider := NewProvider(testAlchemyEndpoint, ProviderAlchemy)

	// process
	params := GetAssetTransfersParams{
		ToAddress:         "0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae",
		ContractAddresses: []string{},
		Category:          []string{"external"},
		Order:             "asc",
		MaxCount:          "0xA",
	}
	result, err := provider.GetAssetTransfers(params)

	// verify
	assert.NoError(t, err)
	assert.NotNil(t, result.PageKey, "page key should not be nil")
	assert.Len(t, result.Transfers, 10, "should contain 10 transfers")

	transfer := result.Transfers[0]
	assert.NotEmpty(t, transfer.Hash, "hash should not be nil")
	assert.NotEmpty(t, transfer.From, "sender should not be nil")
	assert.NotEmpty(t, transfer.To, "receiver should not be nil")
}
