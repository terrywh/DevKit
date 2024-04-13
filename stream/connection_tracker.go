package stream

import (
	"sync"

	"github.com/quic-go/quic-go"
	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/infra/log"
)

type ConnectionTracker interface {
	Enter(conn_id uint64, device_id entity.DeviceID, conn quic.Connection)
	Leave(conn_id uint64, device_id entity.DeviceID, conn quic.Connection)
	Close() error
}

type DefaultConnectionTracker struct {
	mutex *sync.Mutex
	conn  map[uint64]quic.Connection
}

func NewDefaultConnectionTracker() ConnectionTracker {
	return &DefaultConnectionTracker{
		mutex: &sync.Mutex{},
		conn:  make(map[uint64]quic.Connection),
	}
}

func (st *DefaultConnectionTracker) Enter(conn_id uint64, device_id entity.DeviceID, conn quic.Connection) {
	log.Trace("<stream> connection enter: device_id = ", device_id)
	st.mutex.Lock()
	defer st.mutex.Unlock()

	st.conn[conn_id] = conn
}

func (st *DefaultConnectionTracker) Leave(conn_id uint64, device_id entity.DeviceID, conn quic.Connection) {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	delete(st.conn, conn_id)
	log.Trace("<stream> connection leave: device_id = ", device_id)
}

func (st *DefaultConnectionTracker) Close() error {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	for _, conn := range st.conn {
		conn.CloseWithError(quic.ApplicationErrorCode(0), "close")
	}
	return nil
}
