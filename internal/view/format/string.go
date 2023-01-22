package format

import (
	"github.com/dyng/ramen/internal/common"
)

func TruncateText(text string, size int) string {
	if (len(text) > size) {
		return text[:size] + "..."
	} else {
		return text
	}
}

func NormalizeReceiverAddress(receiver *common.Address) string {
	if receiver == nil {
		return "0x"
	} else {
		return receiver.Hex()
	}
}
