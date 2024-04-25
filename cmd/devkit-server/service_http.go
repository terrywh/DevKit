package main

import (
	"context"
	"net/http"
	"time"
)

type HttpService struct {
	mux *http.ServeMux
	svr http.Server
}

func newHttpService() (s *HttpService) {
	s = &HttpService{mux: http.NewServeMux()}
	s.svr = http.Server{Addr: DefaultConfig.Get().Client.Address, Handler: s.mux}
	initHttpServerHandler(s.mux)
	initHttpFileHandler(s.mux)
	return s
}

func (svc *HttpService) Name() string {
	return "http"
}

func (svc *HttpService) Serve(ctx context.Context) {
	go svc.svr.ListenAndServe()

	<-ctx.Done()
	shutdown, cancel := context.WithTimeout(context.Background(), 9*time.Second)
	defer cancel()
	svc.svr.Shutdown(shutdown)
}

func (svc *HttpService) Close() error {
	ctxStop, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	svc.svr.Shutdown(ctxStop) // 10s 超时后，强制停止
	return svc.svr.Close()
}
