package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/stream"
)

type ServiceP2PHandler struct {
	app.StreamHandlerBase
}

func initServiceP2PHandler(mux *stream.ServeMux) {
	handler := &ServiceP2PHandler{}
	mux.HandleFunc("/registry/dial", handler.HandleDial)
}

func (handler *ServiceP2PHandler) HandleDial(ctx context.Context, src *stream.SessionStream) {
	peer := entity.RemotePeer{}
	if err := app.ReadJSON(src.Reader(), &peer); err != nil {
		handler.Respond(src, err)
		return
	}
	log.Println("<ServiceP2PHandler.HandleDial> from: ", peer.DeviceID, peer.Address)
	if !onAuthorize(peer.DeviceID) {
		handler.Respond(src, entity.ErrUnauthorized)
		return
	}
	handler.Respond(src, nil)

	go func(ctx context.Context) {
		data := []byte(peer.DeviceID)
		addr, _ := net.ResolveUDPAddr("udp", peer.Address)
		for i := 0; i < 9; i++ {
			if ctx.Err() != nil {
				break
			}
			log.Println(">>", peer.Address)
			stream.DefaultTransport.WriteTo(data, addr)
			time.Sleep(5 * time.Second)
		}
	}(ctx)
}
