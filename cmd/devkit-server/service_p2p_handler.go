package main

import (
	"context"
	"net"
	"time"

	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/infra/log"
	"github.com/terrywh/devkit/stream"
)

type ServiceP2PHandler struct {
	app.StreamHandlerBase
}

func initServiceP2PHandler(mux *stream.ServeMux) {
	handler := &ServiceP2PHandler{}
	mux.HandleFunc("/relay/dial", handler.HandleDial)
}

func (handler *ServiceP2PHandler) HandleDial(ctx context.Context, src *stream.SessionStream) {
	peer := entity.Server{}
	if err := app.ReadJSON(src.Reader(), &peer); err != nil {
		handler.Respond(src, err)
		return
	}
	log.Info("<devkit-server> reverse dial: ", peer.DeviceID, "(", peer.Address, ")")
	if !onAuthorize(peer.DeviceID) {
		handler.Respond(src, entity.ErrUnauthorized)
		return
	}
	handler.Respond(src, nil)

	go func(ctx context.Context) {
		data := []byte(peer.DeviceID)
		addr, _ := net.ResolveUDPAddr("udp", peer.Address)
		count := 9
		for i := 0; i < count; i++ {
			if ctx.Err() != nil {
				break
			}
			time.Sleep(5 * time.Second)
			log.Trace(">> ", peer.Address, " (", i+1, "/", count, ")")
			stream.DefaultTransport.WriteTo(data, addr)
		}
	}(ctx)
}
