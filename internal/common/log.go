package common

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/log"
)

// Exit prints the message to stderr and exits with status 1
func Exit(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

// PrintMessage prints a message to the console (i.e. stdout)
func PrintMessage(msg string, args ...any) {
	fmt.Printf(msg+"\n", args...)
}

// ErrorStackHandler is a log handler that prints the stack trace of an error
func ErrorStackHandler(h log.Handler) log.Handler {
	return log.FuncHandler(func(r *log.Record) error {
		i := len(r.Ctx) - 1
		if i > 0 {
			e := r.Ctx[i]
			if err, ok := e.(error); ok {
				r.Ctx[i] = fmt.Sprintf("%+v", err)
			}
		}
		return h.Log(r)
	})
}
