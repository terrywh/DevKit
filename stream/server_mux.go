package stream

import (
	"context"
	"encoding/json"

	"github.com/quic-go/quic-go"
	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/infra/log"
)

type StreamHandlerFunc struct {
	fn func(ctx context.Context, ss *SessionStream)
}

func (shf StreamHandlerFunc) ServeStream(ctx context.Context, ss *SessionStream) {
	shf.fn(ctx, ss)
}

type ServeMux struct {
	handler map[string]StreamHandler
}

func NewServeMux() (mux *ServeMux) {
	mux = &ServeMux{
		handler: make(map[string]StreamHandler),
	}
	return mux
}

func (mux ServeMux) Handle(path string, handler StreamHandler) {
	mux.handler[path] = handler
}

func (mux ServeMux) HandleFunc(path string, fn func(ctx context.Context, ss *SessionStream)) {
	mux.handler[path] = StreamHandlerFunc{fn}
}

func (mux ServeMux) ServeStream(ctx context.Context, src *SessionStream) {
	defer src.CloseWrite()
	if ctx.Err() != nil {
		return
	}

	path, err := src.r.ReadString(':')
	if err != nil {
		return
	}
	path = path[:len(path)-1]

	if handler, found := mux.handler[path]; found {
		handler.ServeStream(ctx, src)
	} else {
		log.Warn("<stream> server handler not found: path =", path)
		json.NewEncoder(src).Encode(entity.Response{Error: entity.ErrHandlerNotFound.(*entity.DefaultErrorCode)})
		src.s.CancelRead(quic.StreamErrorCode(entity.ErrSessionNotFound.ErrorCode()))
	}
}
