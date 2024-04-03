package server

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"

	"github.com/quic-go/quic-go"
)

type ServerOptions struct {
	stream.ServerOptions
}

type ServerHandler struct {
	Request stream.Request
	Handler stream.ServerHandler
}

type Server struct {
	options ServerOptions
	handler map[string]ServerHandler
}

func New() *Server {
	svr := &Server{
		handler: make(map[string]ServerHandler),
	}
	svr.handler["StreamFile"] = ServerHandler{
		Request: &stream.StreamFileReq{},
		Handler: &StreamFile{},
	}
	svr.handler["OpenShell"] = ServerHandler{
		Request: &stream.OpenShellReq{},
		Handler: &OpenShell{},
	}
	return svr
}

func (svr *Server) ParseFlags() {
	flag.StringVar(&svr.options.Bind, "bind", "0.0.0.0:12345", "Address server bind to")
	flag.StringVar(&svr.options.Crt, "cert", "./var/cert/server.crt", "Certificate to use for QUIC (server only)")
	flag.StringVar(&svr.options.Key, "pkey", "./var/pkey/server.key", "Private key to use for QUIC (server only)")
}

// func (svr *Server) Handle(command string, handler Handler) {
// 	svr.handler[command] = handler
// }

func (svr *Server) Serve(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	l, err := stream.CreateServer(stream.ServerOptions{})
	if err != nil {
		log.Print(err)
		return
	}
	defer l.Close()

	log.Println("<server.Server> accepting ...")
	for {
		conn, err := l.Accept(ctx)
		if err != nil {
			log.Print(err)
			return
		}
		go svr.serveConn(context.Background(), conn)
	}
}

func (svr *Server) serveConn(ctx context.Context, conn quic.Connection) {
	log.Println("<server.Server> connection started: ", &conn)
SERVING:
	for {
		s, err := conn.AcceptStream(ctx)
		if ae, ok := err.(*quic.ApplicationError); ok && ae.ErrorCode == quic.ApplicationErrorCode(0) {
			break SERVING
		} else if err != nil {
			log.Print("failed to accept stream: ", err)
			return
		}
		go svr.serveStream(context.Background(), s)
	}
	log.Println("<server.Server> connection closed: ", &conn)
}

func (svr *Server) serveStream(ctx context.Context, s quic.Stream) {
	defer s.Close()
	r := bufio.NewReader(s)
	x, err := r.ReadString('\n')
	if err != nil {
		log.Print(err)
		return
	}
	request := strings.SplitN(x, ":", 2)
	if handler, found := svr.handler[request[0]]; found {
		log.Println("<server.Server> handle:", request[0])
		req := reflect.New(reflect.TypeOf(handler.Request).Elem())
		json.Unmarshal([]byte(request[1]), req.Interface())
		handler.Handler.ServeStream(ctx, req.Interface().(stream.Request), r, s)
	} else {
		log.Println("<server.Server> handle not found:", request[0])
		fmt.Fprintln(s, "unknown command")
		s.CancelRead(quic.StreamErrorCode(10000))
	}
}
