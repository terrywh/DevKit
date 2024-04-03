package stream

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"

	"github.com/quic-go/quic-go"
)

type ServerHandler interface {
	ServeStream(ctx context.Context, r *bufio.Reader, w io.Writer)
}

type ServerHandlerFn struct {
	fn func(ctx context.Context, r *bufio.Reader, w io.Writer)
}

func (shf ServerHandlerFn) ServeStream(ctx context.Context, r *bufio.Reader, w io.Writer) {
	shf.fn(ctx, r, w)
}

type ServeMux struct {
	handler map[string]ServerHandler
}

func NewServeMux() (mux *ServeMux) {
	mux = &ServeMux{
		handler: make(map[string]ServerHandler),
	}
	return mux
}

func (mux ServeMux) Handle(path string, handler ServerHandler) {
	mux.handler[path] = handler
}

func (mux ServeMux) HandleFunc(path string, fn func(ctx context.Context, r *bufio.Reader, w io.Writer)) {
	mux.handler[path] = ServerHandlerFn{fn}
}

func (mux ServeMux) ServeStream(ctx context.Context, s quic.Stream) {
	defer s.Close()
	r := bufio.NewReader(s)
	path, err := r.ReadString(':')
	if err != nil {
		log.Print(err)
		return
	}

	if handler, found := mux.handler[path]; found {
		handler.ServeStream(ctx, r, s)
	} else {
		log.Println("<ServeMux.ServeStream> handle not found for path: ", path)
		fmt.Fprintln(s, "unknown command")
		s.CancelRead(quic.StreamErrorCode(10000))
	}
}
