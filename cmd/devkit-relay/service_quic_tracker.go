package main

import (
	"sync"

	"github.com/quic-go/quic-go"
	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/infra/log"
)

func newDiscoveryTracker() *ServiceQuicTracker {
	return &ServiceQuicTracker{
		mutex:  &sync.RWMutex{},
		conn:   make(map[uint64]quic.Connection),
		device: make(map[entity.DeviceID]quic.Connection),
	}
}

type ServiceQuicTracker struct {
	app.StreamHandlerBase
	mutex  *sync.RWMutex
	conn   map[uint64]quic.Connection
	device map[entity.DeviceID]quic.Connection
}

func (dt *ServiceQuicTracker) GetConn(device_id entity.DeviceID) (conn quic.Connection) {
	dt.mutex.RLock()
	defer dt.mutex.RUnlock()

	conn = dt.device[device_id]
	return
}

func (dt *ServiceQuicTracker) Enter(conn_id uint64, device_id entity.DeviceID, conn quic.Connection) {
	log.Trace("<devkit-relay> connection enter: device_id = ", device_id)
	dt.mutex.Lock()
	defer dt.mutex.Unlock()

	dt.conn[conn_id] = conn
	dt.device[device_id] = conn
}

func (dt *ServiceQuicTracker) Leave(conn_id uint64, device_id entity.DeviceID, conn quic.Connection) {
	dt.mutex.Lock()
	defer dt.mutex.Unlock()

	delete(dt.conn, conn_id)
	if dc := dt.device[device_id]; dc == conn {
		delete(dt.device, device_id)
	}
	log.Trace("<devkit-relay> connection leave: device_id = ", device_id)
}

func (dt *ServiceQuicTracker) Close() error {
	dt.mutex.Lock()
	defer dt.mutex.Unlock()
	for _, conn := range dt.conn {
		conn.CloseWithError(quic.ApplicationErrorCode(0), "close")
	}
	// TODO cleanup dt.conn / dt.devs
	return nil
}
