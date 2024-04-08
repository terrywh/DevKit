package stream

import (
	"context"

	"github.com/terrywh/devkit/entity"
)

type Resolver interface {
	Resolve(ctx context.Context, peer *entity.RemotePeer) error
	Serve(ctx context.Context)
	Close() error
}
