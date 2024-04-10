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
	Peer entity.RemotePeer
	Conn quic.Connection
	s    quic.Stream
	r    *bufio.Reader
}

func NewSessionStream(peer *entity.RemotePeer, conn quic.Connection) (ss *SessionStream, err error) {
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

func (ss *SessionStream) RemotePeer() *entity.RemotePeer {
	return &ss.Peer
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

	if _, err = fmt.Fprintf(ss, "%s:", path); err != nil {
		return
	}
	if err = ss.Push(req); err != nil {
		return
	}
	r := entity.HttpResponse{Data: rsp}
	// Decoder 可能读取了更后面的内容，可以使用（响应内容已结束）
	// err = json.NewDecoder(ss).Decode(&r)
	if err = ss.Pull(&r); err != nil {
		return
	}

	if r.Error.Code > 0 {
		err = r.Error
		return
	}
	return nil
}

func (ss *SessionStream) Push(r interface{}) (err error) {
	return json.NewEncoder(ss).Encode(r)
}

func (ss *SessionStream) Pull(r interface{}) (err error) {
	var payload []byte
	if payload, err = ss.r.ReadBytes(byte('\n')); err != nil {
		return
	}
	err = json.Unmarshal(payload, r)
	return
}
