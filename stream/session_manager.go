package stream

import (
	"context"
	"log"
	"sync"

	"github.com/quic-go/quic-go"
	"github.com/terrywh/devkit/entity"
)

type SessionManager struct {
	session  map[entity.DeviceID]*Session
	mutex    *sync.Mutex
	provider ConnectionProvider
}

var DefaultSessionManager *SessionManager

func InitSessionManager(provider ConnectionProvider) (sm *SessionManager) {
	sm = &SessionManager{
		session:  make(map[entity.DeviceID]*Session),
		mutex:    &sync.Mutex{},
		provider: provider,
	}
	DefaultSessionManager = sm
	return
}

func (s *SessionManager) Close() {
	for _, session := range s.session {
		session.Close()
	}
}

func (mgr *SessionManager) Acquire(ctx context.Context, device_id entity.DeviceID) (s *Session, err error) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	var ok bool
	// 复用现有会话
	if s, ok = mgr.session[device_id]; ok {
		return s, nil
	}
	s = &Session{}
	// 建立新会话
	if s.conn, err = mgr.provider.Acquire(ctx, device_id); err != nil {
		return
	}
	mgr.session[device_id] = s
	go func() {
		log.Println("<SessionManager.Acquire> connection: ", &s.conn, " started ...")
		// 监听链接持续时间
		<-s.conn.Context().Done()
		log.Println("<SessionManager.Acquire> connection: ", &s.conn, " closed.")
		mgr.mutex.Lock()
		defer mgr.mutex.Unlock()
		delete(mgr.session, device_id)
		s.Close()
	}()
	return s, nil
}

func (s *SessionManager) AcquireStream(ctx context.Context, device_id entity.DeviceID) (quic.Stream, error) {
	session, err := s.Acquire(ctx, device_id)
	if err != nil {
		return nil, err
	}
	return session.conn.OpenStream()
}
