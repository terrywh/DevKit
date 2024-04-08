package stream

import (
	"log"
	"sync"

	"github.com/quic-go/quic-go"
	"github.com/terrywh/devkit/entity"
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
	log.Println("<ServerTrackerDefault.Enter> connection: ", conn_id, "(device_id = ", device_id, ")")
	st.mutex.Lock()
	defer st.mutex.Unlock()

	st.conn[conn_id] = conn
}

func (st *DefaultConnectionTracker) Leave(conn_id uint64, device_id entity.DeviceID, conn quic.Connection) {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	delete(st.conn, conn_id)
	log.Println("<ServerTrackerDefault.Leave> connection: ", conn_id)
}

func (st *DefaultConnectionTracker) Close() error {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	log.Println("<ServerTrackerDefault.Close>")
	for _, conn := range st.conn {
		conn.CloseWithError(quic.ApplicationErrorCode(0), "close")
	}
	return nil
}
