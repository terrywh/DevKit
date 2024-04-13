package main

import (
	"context"
	"net/http"
	"time"

	"github.com/terrywh/devkit/stream"
)

type HttpService struct {
	mux *http.ServeMux
	svr http.Server
}

func newServiceHttp(mgr stream.SessionManager, mux *stream.ServeMux) (s *HttpService) {
	s = &HttpService{mux: http.NewServeMux()}
	s.svr = http.Server{Addr: DefaultConfig.Get().Client.Address, Handler: s.mux}
	initHttpShellHandler(mgr, s.mux)
	initHttpFileHandler(mgr, s.mux)
	s.mux.Handle("/node_modules/", http.FileServer(http.Dir(".")))
	s.mux.Handle("/", http.FileServer(http.Dir("public")))
	return
}

func (s *HttpService) Serve(ctx context.Context) {
	go s.svr.ListenAndServe()

	<-ctx.Done()
	shutdown, cancel := context.WithTimeout(context.Background(), 9*time.Second)
	defer cancel()
	s.svr.Shutdown(shutdown)
}

func (s *HttpService) Close() error {
	ctxStop, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.svr.Shutdown(ctxStop) // 10s 超时后，强制停止
	return s.svr.Close()
}
