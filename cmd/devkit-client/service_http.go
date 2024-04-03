package main

import (
	"context"
	"net/http"
	"time"

	"github.com/terrywh/devkit/handler"
)

type HttpService struct {
	mux *http.ServeMux
}

func newServiceHttp() (s *HttpService) {
	s = &HttpService{mux: http.NewServeMux()}

	handler.NewClientShellHandler(s.mux)
	s.mux.Handle("/", http.FileServer(http.Dir("www")))
	return
}

func (s *HttpService) Serve(ctx context.Context) {
	mux := http.NewServeMux()

	svr := http.Server{Addr: DefaultOptions.http, Handler: mux}
	go func() {
		<-ctx.Done() // 等待停止通知
		ctxStop, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		svr.Shutdown(ctxStop) // 10s 超时后，强制停止
		svr.Close()
	}()
	svr.ListenAndServe()
}
