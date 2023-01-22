package common

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestAddress_NoError(t *testing.T) {
	addr := common.HexToAddress("0xFABB0ac9d68B0B445fB7357272Ff202C5651694a")
	assert.Equal(t, "0xFABB0ac9d68B0B445fB7357272Ff202C5651694a", addr.Hex())
}
