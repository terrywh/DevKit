package stream

import (
	"bufio"
	"context"
	"io"
)

type ServerHandler interface {
	ServeStream(ctx context.Context, req Request, r *bufio.Reader, w io.Writer)
}

type ClientHandler interface {
	Request() Request
	ServeStream(ctx context.Context, r io.Reader, w io.Writer)
}
