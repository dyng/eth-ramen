package service

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/log"
)

type Contract struct {
	*Account
	abi *abi.ABI
	source string
}

func (c *Contract) GetABI() *abi.ABI {
	return c.abi
}

func (c *Contract) GetSource() string {
	return c.source
}

func (c *Contract) Call(method string, args ...any) ([]any, error) {
	m, ok := c.abi.Methods[method]
	if !ok {
		return nil, fmt.Errorf("Method %s is not found in contract", method)
	}

	if !m.IsConstant() {
		return nil, fmt.Errorf("Method %s is not a constant method", method)
	}

	log.Debug("Try to call contract", "method", method, "args", args)
	return c.service.provider.CallContract(c.address, c.abi, method, args...)
}
