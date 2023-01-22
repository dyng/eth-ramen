package conv

import (
	"encoding/hex"
	"strconv"
)

func HexToInt(s string) (int64, error) {
	return strconv.ParseInt(trim0xPrefix(s), 16, 64)
}

func HexToBytes(s string) ([]byte, error) {
	return hex.DecodeString(trim0xPrefix(s))
}

func has0xPrefix(s string) bool {
	return len(s) >= 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X')
}

func trim0xPrefix(s string) string {
	if has0xPrefix(s) {
		return s[2:]
	} else {
		return s
	}
}
