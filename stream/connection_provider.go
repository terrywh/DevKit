package stream

import (
	"context"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/terrywh/devkit/entity"
)

// ConnectionProvider 链接器（不需要进行多线程保护）
type ConnectionProvider interface {
	Acquire(ctx context.Context, devId entity.DeviceID) (quic.Connection, error)
}

type DirectProvider struct {
	opts DialOptions
	addr map[entity.DeviceID][]string
}

func NewDirectProvider(opts DialOptions) (dp *DirectProvider) {
	if opts.Backoff < time.Second {
		opts.Backoff = 1200 * time.Millisecond
	}
	if opts.Retry == 0 {
		opts.Retry = 3
	}
	dp = &DirectProvider{
		opts: opts,
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
	opts := mgr.opts
	for _, addr := range addrs {
		opts.Address = addr
		conn, err = DefaultTransport.Dial(ctx, opts)
		if err == nil {
			break
		}
	}
	return
}
