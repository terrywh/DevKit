package stream

import (
	"context"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/terrywh/devkit/entity"
)

// ConnectionProvider 链接器（不需要进行多线程保护）
type ConnectionProvider interface {
	Acquire(ctx context.Context, peer *entity.RemotePeer) (quic.Connection, error)
}

type ConnectionProviderCommand interface {
	Execute(ctx context.Context, peer *entity.RemotePeer, conn quic.Connection)
}

type DefaultConnectionProvider struct {
	options DialOptions
	cmd     chan ConnectionProviderCommand
}

func newDefaultConnectionProvider(options *DialOptions) (dp *DefaultConnectionProvider) {
	if options.Address == "" {
		panic("failed to create connection provider: address not provided")
	}
	if options.Backoff < time.Second {
		options.Backoff = 2400 * time.Millisecond
	}
	if options.Retry == 0 {
		options.Retry = 9
	}
	dp = &DefaultConnectionProvider{
		options: *options,
		cmd:     make(chan ConnectionProviderCommand),
	}
	return dp
}

func (provider *DefaultConnectionProvider) Acquire(ctx context.Context, peer *entity.RemotePeer) (conn quic.Connection, err error) {
	opts := provider.options
	opts.Address = peer.Address
	conn, peer.DeviceID, err = DefaultTransport.Dial(ctx, &opts)
	return
}
