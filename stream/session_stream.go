package stream

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"

	"github.com/quic-go/quic-go"
	"github.com/terrywh/devkit/entity"
)

type SessionStream struct {
	peer *entity.RemotePeer
	conn quic.Connection
	s    quic.Stream
	r    *bufio.Reader
}

func NewSessionStream(peer *entity.RemotePeer, conn quic.Connection) (ss *SessionStream, err error) {
	ss = &SessionStream{
		peer: peer,
		conn: conn,
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
	return ss.conn.RemoteAddr()
}

func (ss *SessionStream) RemotePeer() *entity.RemotePeer {
	return ss.peer
}

func (ss *SessionStream) CloseRead() {
	ss.s.CancelRead(quic.StreamErrorCode(0))
}

func (ss *SessionStream) Close() error {
	return ss.s.Close()
}

func (ss *SessionStream) Invoke(ctx context.Context, path string,
	req interface{}, rsp interface{}) (err error) {
	defer ss.Close()

	fmt.Fprintf(ss, "%s:", path)
	if err = json.NewEncoder(ss).Encode(req); err != nil {
		return
	}
	r := entity.HttpResponse{Data: rsp}
	err = json.NewDecoder(ss).Decode(&r)
	if err != nil {
		return
	}
	if r.Error.Code > 0 {
		err = r.Error
		return
	}
	return nil
}
