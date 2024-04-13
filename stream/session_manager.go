package stream

import (
	"context"
	"sync"

	"github.com/quic-go/quic-go"
	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/infra/log"
)

type SessionManager interface {
	EnsureConn(ctx context.Context, peer *entity.Server) (conn quic.Connection, err error)
	Acquire(ctx context.Context, peer *entity.Server) (stream *SessionStream, err error)
	Serve(ctx context.Context)
	Close() error
}

type DefaultSessionManager struct {
	conn     map[entity.DeviceID]quic.Connection
	mutex    *sync.RWMutex
	provider ConnectionProvider
	resolver Resolver
	handler  ConnectionHandler
}

type SessionManagerOptions struct {
	DialOptions
	Resolver Resolver // 默认为空时，不支持 P2P 寻址
	Handler  ConnectionHandler
}

func NewSessionManager(options *SessionManagerOptions) (mgr SessionManager) {
	mgr = &DefaultSessionManager{
		conn:     make(map[entity.DeviceID]quic.Connection),
		mutex:    &sync.RWMutex{},
		provider: newDefaultConnectionProvider(&options.DialOptions),
		resolver: options.Resolver,
		handler:  options.Handler,
	}
	return
}

func (mgr *DefaultSessionManager) Serve(ctx context.Context) {
	mgr.resolver.Serve(ctx)
}

func (mgr *DefaultSessionManager) Close() error {
	mgr.resolver.Close()
	for _, conn := range mgr.conn {
		conn.CloseWithError(quic.ApplicationErrorCode(0), "close")
	}
	return nil
}

func (mgr *DefaultSessionManager) get(peer *entity.Server) (conn quic.Connection) {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()

	if peer.DeviceID != "" {
		conn = mgr.conn[peer.DeviceID]
	}
	return conn
}
func (mgr *DefaultSessionManager) put(peer *entity.Server, conn quic.Connection) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	mgr.conn[peer.DeviceID] = conn
}

func (mgr *DefaultSessionManager) del(peer *entity.Server, conn quic.Connection) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	c := mgr.conn[peer.DeviceID]
	if c == conn {
		delete(mgr.conn, peer.DeviceID)
	}
}

func (mgr *DefaultSessionManager) EnsureConn(ctx context.Context, peer *entity.Server) (conn quic.Connection, err error) {
	if conn = mgr.get(peer); conn != nil {
		return
	}
	if peer.Address == "" {
		if err = mgr.resolver.Resolve(ctx, peer); err != nil {
			return
		}
		log.Debug("<stream> device address resolved: ", peer.Address)
	}
	// 建立新会话
	if conn, err = mgr.provider.Acquire(ctx, peer); err != nil {
		return
	}
	mgr.put(peer, conn)
	go func(conn quic.Connection, peer entity.Server, handler ConnectionHandler) {
		if handler != nil {
			handler.ServeConn(context.Background(), conn)
		} else {
			<-conn.Context().Done() // 监听链接持续时间
		}
		conn.CloseWithError(quic.ApplicationErrorCode(0), "close")
		mgr.del(&peer, conn)
	}(conn, *peer, mgr.handler)
	return conn, nil
}

func (mgr *DefaultSessionManager) Acquire(ctx context.Context, peer *entity.Server) (ss *SessionStream, err error) {
	var conn quic.Connection
	conn, err = mgr.EnsureConn(ctx, peer)
	if err != nil {
		return nil, err
	}
	return NewSessionStream(peer, conn)
}
