package stream

import (
	"bufio"
	"net"

	"github.com/quic-go/quic-go"
	"github.com/terrywh/devkit/entity"
)

type SessionStream struct {
	Peer entity.Server
	Conn quic.Connection
	s    quic.Stream
	r    *bufio.Reader
}

func NewSessionStream(peer *entity.Server, conn quic.Connection) (ss *SessionStream, err error) {
	ss = &SessionStream{
		Peer: *peer,
		Conn: conn,
	}
	if ss.s, err = conn.OpenStream(); err != nil {
		return
	}
	ss.r = bufio.NewReader(ss.s)
	return
}

func (ss *SessionStream) Reader() *bufio.Reader {
	return ss.r
}

func (ss *SessionStream) Read(data []byte) (int, error) {
	return ss.r.Read(data)
}

func (ss *SessionStream) Write(data []byte) (int, error) {
	return ss.s.Write(data)
}

func (ss *SessionStream) RemoteAddr() net.Addr {
	return ss.Conn.RemoteAddr()
}

func (ss *SessionStream) RemotePeer() *entity.Server {
	return &ss.Peer
}

func (ss *SessionStream) CloseRead() {
	ss.s.CancelRead(quic.StreamErrorCode(0))
}

func (ss *SessionStream) CloseWrite() error {
	return ss.s.Close()
}
