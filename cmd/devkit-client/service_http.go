package main

import (
	"context"
	"net/http"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/terrywh/devkit/stream"
)

type HttpService struct {
	mux *http.ServeMux
	svr http.Server
}

type SessionManager interface {
	EnsureConn(conn quic.Connection, err error)
}

func newServiceHttp(mgr stream.SessionManager) (s *HttpService) {
	s = &HttpService{mux: http.NewServeMux()}
	s.svr = http.Server{Addr: DefaultConfig.Get().Client.Address, Handler: s.mux}
	newClientShellHandler(mgr, s.mux)
	s.mux.Handle("/node_modules/", http.FileServer(http.Dir(".")))
	s.mux.Handle("/", http.FileServer(http.Dir("www")))
	return
}

func (s *HttpService) Serve(ctx context.Context) {
	s.svr.ListenAndServe()
}

func (s *HttpService) Close() error {
	ctxStop, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.svr.Shutdown(ctxStop) // 10s 超时后，强制停止
	return s.svr.Close()
}
