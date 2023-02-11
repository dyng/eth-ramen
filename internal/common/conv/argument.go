package conv

import (
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

// PackArgument packs a single argument into a string
func PackArgument(t abi.Type, v any) (string, error) {
	switch t.T {
	case abi.StringTy:
		if v, ok := v.(string); ok {
			return v, nil
		} else {
			return "", errors.Errorf("cannot convert %v to string", v)
		}
	case abi.IntTy, abi.UintTy:
		return formatInteger(t, v)
	case abi.BoolTy:
		if v, ok := v.(bool); ok {
			return strconv.FormatBool(v), nil
		} else {
			return "", errors.Errorf("cannot convert %v to bool", v)
		}
	case abi.AddressTy:
		if v, ok := v.(common.Address); ok {
			return v.Hex(), nil
		} else {
			return "", errors.Errorf("cannot convert %v to address", v)
		}
	case abi.HashTy:
		if v, ok := v.(common.Hash); ok {
			return v.Hex(), nil
		} else {
			return "", errors.Errorf("cannot convert %v to hash", v)
		}
	default:
		return "", errors.Errorf("unsupported argument type %v", t.T)
	}
}

// UnpackArgument converts string format of a value into the Go type corresponding to given argument type.
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
		return nil, errors.Errorf("unsupported argument type %v", t.T)
	}
}

func formatInteger(t abi.Type, v any) (string, error) {
	if t.T == abi.UintTy {
		switch t.Size {
		case 8, 16, 32, 64:
			if v, ok := v.(uint64); ok {
				return strconv.FormatUint(v, 10), nil
			} else {
				return "", errors.Errorf("cannot convert %v to uint64", v)
			}
		default:
			if v, ok := v.(*big.Int); ok {
				return v.String(), nil
			} else {
				return "", errors.Errorf("cannot convert %v to *big.Int", v)
			}
		}
	} else {
		switch t.Size {
		case 8, 16, 32, 64:
			if v, ok := v.(int64); ok {
				return strconv.FormatInt(v, 10), nil
			} else {
				return "", errors.Errorf("cannot convert %v to int64", v)
			}
		default:
			if v, ok := v.(*big.Int); ok {
				return v.String(), nil
			} else {
				return "", errors.Errorf("cannot convert %v to *big.Int", v)
			}
		}
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
				return nil, errors.Errorf("cannot parse %s as type %v", s, t.T)
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
				return nil, errors.Errorf("cannot parse %s as type %v", s, t.T)
			} else {
				return i, nil
			}
		}
	}
}
