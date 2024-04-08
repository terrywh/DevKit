package stream

import (
	"bufio"
	"context"
	"log"
	"sync/atomic"

	"github.com/quic-go/quic-go"
	"github.com/terrywh/devkit/entity"
)

type ConnectionHandler interface {
	ServeConn(ctx context.Context, conn quic.Connection)
	Close() error
}

type StreamHandler interface {
	ServeStream(ctx context.Context, stream *SessionStream)
}

type DefaultConnectionHandler struct {
	Tracker ConnectionTracker
	Handler StreamHandler

	connid atomic.Uint64
}

func (svr *DefaultConnectionHandler) enter(conn quic.Connection) (conn_id uint64, device_id entity.DeviceID) {
	certs := conn.ConnectionState().TLS.PeerCertificates
	device_id = DeviceIDFromCert(certs[0].Raw)
	conn_id = svr.connid.Add(1)
	svr.Tracker.Enter(conn_id, device_id, conn)
	return
}

func (svr *DefaultConnectionHandler) leave(conn_id uint64, device_id entity.DeviceID, conn quic.Connection) {
	svr.Tracker.Leave(conn_id, device_id, conn)
}

func (svr *DefaultConnectionHandler) ServeConn(ctx context.Context, conn quic.Connection) {
	conn_id, device_id := svr.enter(conn)
SERVING:
	for {
		s, err := conn.AcceptStream(ctx)
		if ae, ok := err.(*quic.ApplicationError); ok && ae.ErrorCode == quic.ApplicationErrorCode(0) {
			break SERVING
		} else if err != nil {
			if _, ok := err.(*quic.ApplicationError); !ok {
				log.Print("failed to accept stream: ", err)
			}
			return
		}

		go svr.Handler.ServeStream(context.Background(), &SessionStream{
			peer: &entity.RemotePeer{DeviceID: device_id, Address: conn.RemoteAddr().String()},
			conn: conn, s: s, r: bufio.NewReader(s),
		})
	}

	svr.leave(conn_id, device_id, conn)
}

func (svr *DefaultConnectionHandler) Close() error {
	return svr.Tracker.Close()
}
