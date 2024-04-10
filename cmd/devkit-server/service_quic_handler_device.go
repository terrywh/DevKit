package main

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/entity"
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
	handler.major, handler.minor, handler.build = handler.initDeviceVersion()
	mux.HandleFunc("/device/query", handler.HandleQuery)
	// TODO cleanup
	return handler
}

func (hss *DeviceHandler) HandleQuery(ctx context.Context, ss *stream.SessionStream) {
	log.Println("<SystemHandler.HandleQuery> device =", ss.RemotePeer().DeviceID)
	hss.Respond(ss, entity.RemotePeer{
		DeviceID: DefaultConfig.Get().DeviceID(),
		System:   runtime.GOOS,
		Arch:     runtime.GOARCH,
		Version:  fmt.Sprintf("%d.%d.%d", hss.major, hss.minor, hss.build),
	})
}
