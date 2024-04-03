package main

import (
	"context"
	"log"

	"github.com/quic-go/quic-go"
	"github.com/terrywh/devkit/handler"
	"github.com/terrywh/devkit/stream"
)

type QuicService struct {
	mux *stream.ServeMux
}

func newQuicService() (qs *QuicService) {
	qs = &QuicService{mux: stream.NewServeMux()}
	handler.NewServerShellHandler(qs.mux)
	return
}

func (svc *QuicService) Serve(ctx context.Context) {
	listener, err := stream.DefaultTransport.CreateServer(stream.ServerOptions{
		Certificate:         DefaultOptions.cert,
		PrivateKey:          DefaultOptions.pkey,
		ApplicationProtocol: "devkit",
	})
	if err != nil {
		panic("failed to create server (QUIC): " + err.Error())
	}
	defer listener.Close()
	log.Println("<QuicService.Serve> accepting... ")
SERVING:
	for {
		conn, err := listener.Accept(ctx)
		if err != nil {
			log.Print(err)
			break SERVING
		}
		go svc.ServeConn(context.Background(), conn)
	}
	log.Println("<QuicService.Serve> closed.")
}
func (svc *QuicService) ServeConn(ctx context.Context, conn quic.Connection) {
	log.Println("<QuicService.ServeConn> connection started: ", &conn)
SERVING:
	for {
		s, err := conn.AcceptStream(ctx)
		if ae, ok := err.(*quic.ApplicationError); ok && ae.ErrorCode == quic.ApplicationErrorCode(0) {
			break SERVING
		} else if err != nil {
			log.Print("failed to accept stream: ", err)
			return
		}
		go svc.ServeStream(context.Background(), s)
	}
	log.Println("<QuicService.ServeConn> connection closed: ", &conn)
}

func (svc *QuicService) ServeStream(ctx context.Context, s quic.Stream) {
	svc.mux.ServeStream(ctx, s)
}
