package format

import (
	"fmt"

	"github.com/dyng/ramen/internal/common"
)

func TruncateText(text string, size int) string {
	if len(text) > size {
		return text[:size] + "..."
	} else {
		return text
	}
}

func NormalizeReceiverAddress(receiver *common.Address) string {
	if receiver == nil {
		return "0x0"
	} else {
		return receiver.Hex()
	}
}

func FineErrorMessage(msg string, args ...any) string {
	if len(args) == 0 {
		return msg
	}

	message := ""
	last := len(args) - 1
	err, ok := args[last].(error)
	if ok {
		message = fmt.Sprintf(msg, args[:last]...)
		message += fmt.Sprintf("\n\nError:\n%s", err)
		message += "\n\nPlease check the log files for more details."
	} else {
		message = fmt.Sprintf(msg, args...)
	}

	return message
}
