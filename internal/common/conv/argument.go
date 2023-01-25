package conv

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

func UnpackArgument(t abi.Type, s string) (any, error) {
	switch t.T {
	case abi.StringTy:
		return s, nil
	case abi.IntTy, abi.UintTy:
		return parseInteger(t, s)
	case abi.BoolTy:
		return strconv.ParseBool(s)
	case abi.AddressTy:
		return common.HexToAddress(s), nil
	case abi.HashTy:
		return common.HexToHash(s), nil
	default:
		return nil, fmt.Errorf("unsupported argument type %v", t.T)
	}
}

func parseInteger(t abi.Type, s string) (any, error) {
	if t.T == abi.UintTy {
		switch t.Size {
		case 8, 16, 32, 64:
			return strconv.ParseUint(s, 10, 64)
		default:
			i, ok := new(big.Int).SetString(s, 10)
			if !ok {
				return nil, fmt.Errorf("cannot parse %s as type %v", s, t.T)
			} else {
				return i, nil
			}
		}
	} else {
		switch t.Size {
		case 8, 16, 32, 64:
			return strconv.ParseInt(s, 10, 64)
		default:
			i, ok := new(big.Int).SetString(s, 10)
			if !ok {
				return nil, fmt.Errorf("cannot parse %s as type %v", s, t.T)
			} else {
				return i, nil
			}
		}
	}
}
