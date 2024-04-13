package main

import (
	"context"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/infra/log"
	"github.com/terrywh/devkit/stream"
)

type Resolver struct {
	options stream.DialOptions
	command chan ResolverCommand
}

type ResolverCommand interface {
	Execute(ctx context.Context, peer *entity.Server, conn quic.Connection)
}

func newResolver(options *stream.DialOptions) (r *Resolver) {
	r = &Resolver{
		command: make(chan ResolverCommand),
		options: *options,
	}
	r.options.Retry = 3
	r.options.Backoff = 1200 * time.Millisecond
	return
}

func (r *Resolver) Serve(ctx context.Context) {
	opts := r.options
	var err error
	var conn quic.Connection
	peer := &entity.Server{
		Address: r.options.Address,
	}
	log.Trace("<devkit-client> resolver started ...")
SERVING:
	for {
		conn, peer.DeviceID, err = stream.DefaultTransport.Dial(ctx, &opts)
		if ctx.Err() != nil {
			break SERVING
		}
		if err != nil {
			continue
		}
		// 追踪连接或重连
	EXECUTING:
		for {
			select {
			case <-conn.Context().Done():
				break EXECUTING
			case <-ctx.Done():
				break SERVING
			case cmd := <-r.command:
				if cmd == nil {
					break SERVING
				}
				cmd.Execute(ctx, peer, conn)
			}
		}
	}
	log.Trace("<devkit-client> resolver closed.")
}

func (r *Resolver) Resolve(ctx context.Context, peer *entity.Server) error {
	cmd := newResolverCommandDial(peer)
	r.command <- cmd
	reply := <-cmd.C
	// 整段均为同一个指针，其值已经填充了
	// peer.Address = reply.P.Address
	return reply.E
}

func (r *Resolver) Close() error {
	close(r.command)
	return nil
}

type ResolverCommandDialReply struct {
	E error
	P *entity.Server
}

type ResolverCommandDial struct {
	C chan *ResolverCommandDialReply
	P *entity.Server
}

func newResolverCommandDial(peer *entity.Server) *ResolverCommandDial {
	return &ResolverCommandDial{
		C: make(chan *ResolverCommandDialReply),
		P: peer,
	}
}

func (rcd *ResolverCommandDial) Execute(ctx context.Context, peer *entity.Server, conn quic.Connection) {
	ss, err := stream.NewSessionStream(peer, conn)

	r := &ResolverCommandDialReply{
		P: rcd.P, // 引用原指针，Invoke 回填
	}
	if err != nil {
		r.E = err
		rcd.C <- r
		return
	}
	r.E = app.Invoke(ctx, ss, "/relay/dial", rcd.P, r.P)
	rcd.C <- r
}
