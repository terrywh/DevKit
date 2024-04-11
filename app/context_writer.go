package app

import "context"

type ContextDiscardWriter struct {
	Ctx context.Context
}

func (w ContextDiscardWriter) Write(payload []byte) (size int, err error) {
	size = len(payload)
	err = w.Ctx.Err()
	return
}
