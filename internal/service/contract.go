package service

import (
	"strings"

	"github.com/dyng/ramen/internal/common"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/log"
	"github.com/pkg/errors"
)

// Contract represents a smart contract deployed on Ethereum network.
type Contract struct {
	*Account
	abi *abi.ABI
	source string
}

// HasABI returns true if this contract has a known ABI.
func (c *Contract) HasABI() bool {
	return c.abi != nil
}

// GetABI returns ABI of this contract, may be nil if ABI is unknown.
func (c *Contract) GetABI() *abi.ABI {
	return c.abi
}

// GetSource returns source of this contract, may be empty if source cannot be retrieved.
func (c *Contract) GetSource() string {
	return c.source
}

// ImportABI generates ABI from a json representation of ABI.
func (c *Contract) ImportABI(abiJson string) error {
	log.Debug("Try to parse abi json", "json", abiJson)
	parsedAbi, err := abi.JSON(strings.NewReader(abiJson))
	if err != nil {
		return errors.WithStack(err)
	}

	c.abi = &parsedAbi

	return nil
}

// Call invokes a constant method of this contract. The arguments should be unpacked into correct type.
func (c *Contract) Call(method string, args ...any) ([]any, error) {
	m, ok := c.abi.Methods[method]
	if !ok {
		return nil, errors.Errorf("Method %s is not found in contract", method)
	}

	if !m.IsConstant() {
		return nil, errors.Errorf("Method %s is not a constant method", method)
	}

	log.Debug("Try to call contract", "method", method, "args", args)
	return c.service.provider.CallContract(c.address, c.abi, method, args...)
}

// Send invokes a non-constant method of this contract. This method will sign and send the transaction to the network.
func (c *Contract) Send(signer *Signer, method string, args ...any) (common.Hash, error) {
	_, ok := c.abi.Methods[method]
	if !ok {
		return common.Hash{}, errors.Errorf("Method %s is not found in contract", method)
	}

	return signer.CallContract(c.GetAddress(), c.abi, method, args...)
}
