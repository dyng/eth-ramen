package service

import (
	"crypto/ecdsa"

	"github.com/dyng/ramen/internal/common"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/params"
	"github.com/pkg/errors"
)

type Signer struct {
	*Account
	PrivateKey *ecdsa.PrivateKey
}

func (s *Signer) TransferTo(address common.Address, amount common.BigInt) (common.Hash, error) {
	gasPrice, err := s.service.provider.GetGasPrice()
	if err != nil {
		return common.Hash{}, err
	}

	txnReq := &common.TxnRequest{
		PrivateKey: s.PrivateKey,
		To:         &address,
		Value:      amount,
		GasLimit:   params.TxGas,
		GasPrice:   gasPrice,
	}

	return s.service.provider.SendTransaction(txnReq)
}

func (s *Signer) CallContract(address common.Address, abi *abi.ABI, method string, args ...any) (common.Hash, error) {
	gasPrice, err := s.service.provider.GetGasPrice()
	if err != nil {
		return common.Hash{}, err
	}

	input, err := abi.Pack(method, args...)
	if err != nil {
		return common.Hash{}, errors.WithStack(err)
	}

	gasLimit, err := s.service.provider.EstimateGas(address, s.address, input)
	if err != nil {
		return common.Hash{}, err
	}

	txnReq := &common.TxnRequest{
		PrivateKey: s.PrivateKey,
		To:         &address,
		GasLimit:   gasLimit,
		GasPrice:   gasPrice,
		Data:       input,
	}

	return s.service.provider.SendTransaction(txnReq)
}
