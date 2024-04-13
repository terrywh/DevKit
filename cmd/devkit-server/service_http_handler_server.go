package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/infra"
)

type HttpServerHandler struct {
	app.HttpHandlerBase
	major, minor, build uint32
}

func initHttpServerHandler(mux *http.ServeMux) *HttpServerHandler {
	handler := &HttpServerHandler{}
	handler.major, handler.minor, handler.build = infra.DefaultSystem.Version()
	mux.HandleFunc("/server/query", handler.HandleQuery)
	return handler
}

func (handler *HttpServerHandler) HandleQuery(w http.ResponseWriter, r *http.Request) {
	handler.Respond(w, entity.Server{
		DeviceID: DefaultConfig.Get().DeviceID(),
		Pid:      os.Getpid(),
		System:   runtime.GOOS,
		Arch:     runtime.GOARCH,
		Version:  fmt.Sprintf("%d.%d.%d", handler.major, handler.minor, handler.build),
	})
}
