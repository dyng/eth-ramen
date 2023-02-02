package service

import (
	"testing"

	"github.com/dyng/ramen/internal/common/conv"
	"github.com/stretchr/testify/assert"
)

func TestTransfer(t *testing.T) {
	// prepare
	serv := NewTestService()

	// sender
	senderKey := "0xde9be858da4a475276426320d5e9262ecfc3ba460bfac56360bfa6c4c28b4ee0"
	sender, _ := serv.GetSigner(senderKey)
	balance, _ := sender.GetBalanceForce()
	assert.EqualValues(t, conv.FromEther(10000), balance, "sender should have 10000 eth")

	// receiver
	receiver, _ := serv.GetAccount("0x8626f6940E2eb28930eFb4CeF49B2d1F2C9C1199")
	balance, _ = receiver.GetBalanceForce()
	assert.EqualValues(t, conv.FromEther(10000), balance, "receiver should have 10000 eth")

	// process
	_, err := sender.TransferTo(receiver.GetAddress(), conv.FromEther(10))

	// verify
	assert.NoError(t, err)
	balance, _ = sender.GetBalanceForce()
	assert.LessOrEqual(t, balance.Cmp(conv.FromEther(9990)), 0, "sender should have less than 9990 eth")
	balance, _ = receiver.GetBalanceForce()
	assert.EqualValues(t, conv.FromEther(10010), balance, "sender should have 10010 eth")
}
