package stream

import (
	"context"
	"log"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/terrywh/devkit/entity"
)

// ConnectionProvider 链接器（不需要进行多线程保护）
type ConnectionProvider interface {
	Acquire(ctx context.Context, peer *entity.RemotePeer) (quic.Connection, error)
	Serve(ctx context.Context)
}

type ConnectionProviderCommand interface {
	Execute(ctx context.Context, conn quic.Connection)
}

type ConnectionProviderDialResult struct {
	E error
	P *entity.RemotePeer
}

type ConnectionProviderDial struct {
	P *entity.RemotePeer
	C chan ConnectionProviderDialResult
}

func (cpd ConnectionProviderDial) Execute(ctx context.Context, conn quic.Connection) {
	ss, err := NewSessionStream(&entity.RemotePeer{
		Address: "127.0.0.1:18082",
	}, conn)

	r := ConnectionProviderDialResult{}
	if err != nil {
		r.E = err
		cpd.C <- r
		return
	}
	r.E = ss.Invoke(ctx, "/registry/relay/dial", cpd.P, r.P)
	cpd.C <- r
}

type DefaultConnectionProvider struct {
	options *DialOptions
	conn    quic.Connection
	dial    chan ConnectionProviderCommand
}

func NewDefaultConnectionProvider(options *DialOptions) (dp *DefaultConnectionProvider) {
	if options.Address == "" {
		options.Address = "127.0.0.1:18082" // TODO 公共 REGISTRY 服务
	}
	if options.Backoff < time.Second {
		options.Backoff = 2400 * time.Millisecond
	}
	if options.Retry == 0 {
		options.Retry = 9
	}
	dp = &DefaultConnectionProvider{options: options}
	go dp.Serve(context.Background())
	return dp
}

func (mgr *DefaultConnectionProvider) Serve(ctx context.Context) {
	opts := *mgr.options
	opts.Retry = 3
	var err error
SERVING:
	for {
		mgr.conn, _, err = DefaultTransport.Dial(ctx, opts)
		if err != nil {
			continue
		}
		// 追踪连接或重连
		select {
		case <-mgr.conn.Context().Done():
			time.Sleep(5 * time.Second)
			continue
		case <-ctx.Done():
			break SERVING
		case peer := <-mgr.dial:
			// TODO peer query address
			log.Println(peer)
		}
	}
}

func (mgr *DefaultConnectionProvider) Acquire(ctx context.Context, peer *entity.RemotePeer) (conn quic.Connection, err error) {
	if peer.Address == "" { // 通过 registry 补充 server 地址（若未指定）
		cpd := ConnectionProviderDial{
			P: &entity.RemotePeer{
				DeviceID: peer.DeviceID,
			},
			C: make(chan ConnectionProviderDialResult),
		}
		mgr.dial <- cpd
		rst := <-cpd.C
		if rst.E != nil {
			err = rst.E
			return
		}
		if rst.P.DeviceID != entity.DeviceID(peer.Address) {
			err = entity.ErrSessionNotFound
			return
		}
		peer.Address = rst.P.Address
	}

	opts := *mgr.options
	opts.Address = peer.Address
	conn, peer.DeviceID, err = DefaultTransport.Dial(ctx, opts)
	return
}
