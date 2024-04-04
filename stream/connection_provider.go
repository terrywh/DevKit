package stream

import (
	"context"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/terrywh/devkit/entity"
)

type ConnectionProvider interface {
	Acquire(ctx context.Context, devId entity.DeviceID) (quic.Connection, error)
}

type DirectProvider struct {
	addr map[entity.DeviceID][]string
}

func NewDirectProvider() (dp *DirectProvider) {
	dp = &DirectProvider{
		addr: make(map[entity.DeviceID][]string),
	}
	return dp
}

func (mgr *DirectProvider) Acquire(ctx context.Context, id entity.DeviceID) (conn quic.Connection, err error) {
	var ok bool
	var addr []string

	if addr, ok = mgr.addr[id]; !ok {
		addr = mgr.acquireAddress(ctx, id)
		mgr.addr[id] = addr
	}
	conn, err = mgr.dial(ctx, addr)
	return
}

func (mgr *DirectProvider) acquireAddress(_ context.Context, device_id entity.DeviceID) []string {
	return []string{string(device_id)} // host:port 使用目标地址做标识，直接链接
}

func (mgr *DirectProvider) dial(ctx context.Context, addrs []string) (conn quic.Connection, err error) {
	for _, addr := range addrs {
		conn, err = DefaultTransport.DialEx(ctx, DialOptions{Address: addr, Retry: 2, Backoff: 1200 * time.Millisecond})
		if err == nil {
			break
		}
	}
	return
}
