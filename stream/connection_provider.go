package stream

import (
	"context"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/terrywh/devkit/entity"
)

// ConnectionProvider 链接器（不需要进行多线程保护）
type ConnectionProvider interface {
	Acquire(ctx context.Context, peer *entity.Server) (quic.Connection, error)
}

type DefaultConnectionProvider struct {
	options DialOptions
}

func newDefaultConnectionProvider(options *DialOptions) (dp *DefaultConnectionProvider) {
	if options.Address == "" {
		panic("failed to create connection provider: address not provided")
	}
	dp = &DefaultConnectionProvider{
		options: *options,
	}
	if dp.options.Backoff < time.Second {
		dp.options.Backoff = 2400 * time.Millisecond
	}
	if dp.options.Retry == 0 {
		dp.options.Retry = 9
	}
	return dp
}

func (provider *DefaultConnectionProvider) Acquire(ctx context.Context, peer *entity.Server) (conn quic.Connection, err error) {
	opts := provider.options
	opts.Address = peer.Address
	conn, peer.DeviceID, err = DefaultTransport.Dial(ctx, &opts)
	return
}
