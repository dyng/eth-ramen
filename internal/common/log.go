package common

import (
	"fmt"

	"github.com/ethereum/go-ethereum/log"
)

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
