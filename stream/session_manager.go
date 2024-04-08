package stream

import (
	"context"
	"log"
	"sync"

	"github.com/quic-go/quic-go"
	"github.com/terrywh/devkit/entity"
)

type SessionManager interface {
	EnsureConn(ctx context.Context, peer *entity.RemotePeer) (conn quic.Connection, err error)
	Acquire(ctx context.Context, peer *entity.RemotePeer) (stream *SessionStream, err error)
	Serve(ctx context.Context)
	Close() error
}

type DefaultSessionManager struct {
	conn     map[entity.DeviceID]quic.Connection
	mutex    *sync.Mutex
	provider ConnectionProvider
}

func NewSessionManager(provider ConnectionProvider) (mgr SessionManager) {
	mgr = &DefaultSessionManager{
		conn:     make(map[entity.DeviceID]quic.Connection),
		mutex:    &sync.Mutex{},
		provider: provider,
	}
	return
}

func (mgr *DefaultSessionManager) Serve(ctx context.Context) {
	mgr.provider.Serve(ctx)
}

func (mgr *DefaultSessionManager) Close() error {
	for _, conn := range mgr.conn {
		conn.CloseWithError(quic.ApplicationErrorCode(0), "close")
	}
	return nil
}

func (mgr *DefaultSessionManager) EnsureConn(ctx context.Context, peer *entity.RemotePeer) (conn quic.Connection, err error) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	var ok bool
	// 复用现有会话
	if conn, ok = mgr.conn[peer.DeviceID]; peer.DeviceID != "" && ok {
		return conn, nil
	}
	// 建立新会话
	if conn, err = mgr.provider.Acquire(ctx, peer); err != nil {
		return
	}
	mgr.conn[peer.DeviceID] = conn
	go func() {
		ctx := conn.Context()
		log.Println("<SessionManager.Acquire> connection: ", &conn, " started ...")
		// 监听链接持续时间
		<-ctx.Done()
		log.Println("<SessionManager.Acquire> connection: ", &conn, " closed.")

		mgr.mutex.Lock()
		defer mgr.mutex.Unlock()
		delete(mgr.conn, peer.DeviceID)
		conn.CloseWithError(quic.ApplicationErrorCode(0), "close")
	}()
	return conn, nil
}

func (mgr *DefaultSessionManager) Acquire(ctx context.Context, peer *entity.RemotePeer) (ss *SessionStream, err error) {
	var conn quic.Connection
	conn, err = mgr.EnsureConn(ctx, peer)
	if err != nil {
		return nil, err
	}
	return NewSessionStream(peer, conn)
}
