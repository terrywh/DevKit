package stream

import (
	"github.com/quic-go/quic-go"
)

type Session struct {
	conn quic.Connection
}

func (s *Session) Close() error {
	return s.conn.CloseWithError(quic.ApplicationErrorCode(0), "close")
}

func (s *Session) OpenStream() (quic.Stream, error) {
	return s.conn.OpenStream()
}
