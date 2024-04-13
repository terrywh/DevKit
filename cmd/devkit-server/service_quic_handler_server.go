package main

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/infra"
	"github.com/terrywh/devkit/stream"
)

type DeviceHandler struct {
	app.StreamHandlerBase
	major uint32
	minor uint32
	build uint32
}

func initDeviceHandler(mux *stream.ServeMux) *DeviceHandler {
	handler := &DeviceHandler{}
	handler.major, handler.minor, handler.build = infra.DefaultSystem.Version()
	mux.HandleFunc("/server/query", handler.HandleQuery)
	// TODO cleanup
	return handler
}

func (hss *DeviceHandler) HandleQuery(ctx context.Context, ss *stream.SessionStream) {
	hss.Respond(ss, entity.Server{
		DeviceID: DefaultConfig.Get().DeviceID(),
		Pid:      os.Getpid(),
		System:   runtime.GOOS,
		Arch:     runtime.GOARCH,
		Version:  fmt.Sprintf("%d.%d.%d", hss.major, hss.minor, hss.build),
	})
}
