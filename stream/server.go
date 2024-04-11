package stream

import (
	"context"
	"log"

	"github.com/quic-go/quic-go"
	"github.com/terrywh/devkit/entity"
)

type Server struct {
	listener *quic.Listener
	handler  ConnectionHandler
}

type ServerOptions struct {
	Handler     ConnectionHandler
	Authorize   func(device_id entity.DeviceID) bool
	Certificate string
	PrivateKey  string

	ApplicationProtocol string
}

func NewServer(options *ServerOptions) (svr *Server, err error) {
	if options.Handler == nil {
		err = entity.ErrInvalidArguments
		return
	}
	svr = &Server{handler: options.Handler}
	svr.listener, err = DefaultTransport.createListener(options)
	return
}

func (svr *Server) Serve(ctx context.Context) {
	defer svr.listener.Close()
	log.Println("<Server.Serve> accepting... ")
SERVING:
	for {
		conn, err := svr.listener.Accept(ctx)
		if err != nil {
			break SERVING
		}
		go svr.handler.ServeConn(ctx, conn)
	}
	log.Println("<Server.Serve> closed.")
}

func (svr *Server) Close() error {
	svr.handler.Close()
	return svr.listener.Close()
}
