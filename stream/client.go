package stream

import (
	"context"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/terrywh/devkit/infra/log"
)

type Client struct {
	handler ConnectionHandler
	options DialOptions
}

type ClientOptions struct {
	Handler ConnectionHandler
	DialOptions
}

func NewClient(options *ClientOptions) (cli *Client) {
	cli = &Client{
		handler: options.Handler,
		options: options.DialOptions,
	}
	return cli
}

func (cli *Client) Serve(ctx context.Context) {
	var conn quic.Connection
	// var device_id entity.DeviceID
	var err error
SERVING:
	for {
		if ctx.Err() != nil {
			break SERVING
		}
		conn /* device_id */, _, err = DefaultTransport.Dial(ctx, &DialOptions{
			Address:     cli.options.Address, // TODO 公共 REGISTRY 服务
			Certificate: cli.options.Certificate,
			PrivateKey:  cli.options.PrivateKey,
			Retry:       3,
			Backoff:     1200 * time.Millisecond,
		})
		if err != nil {
			log.Warn("<stream> failed to dial relay: ", err)
			continue
		}
		cli.handler.ServeConn(ctx, conn)
		conn.CloseWithError(quic.ApplicationErrorCode(0), "")
	}

}

func (cli *Client) Close() error {
	return cli.handler.Close()
}
