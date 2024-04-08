package main

import (
	"context"
	"encoding/json"
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

func (handler *ServiceP2PHandler) HandleDial(ctx context.Context, ss *stream.SessionStream) {
	peer := entity.RemotePeer{}
	if err := json.NewDecoder(ss).Decode(&peer); err != nil {
		handler.Respond(ss, err)
		return
	}
	log.Println("<ServiceP2PHandler.HandleDial> dial: ", peer.DeviceID, " from: ", ss.RemotePeer().DeviceID)
	if !onAuthorize(peer.DeviceID) {
		handler.Respond(ss, entity.ErrUnauthorized)
		return
	}
	handler.Respond(ss, nil)

	go func(ctx context.Context) {
		data := []byte(peer.DeviceID)
		addr, _ := net.ResolveUDPAddr("udp", peer.Address)
		for i := 0; i < 9; i++ {
			if ctx.Err() != nil {
				break
			}
			time.Sleep(3 * time.Second)
			stream.DefaultTransport.WriteTo(data, addr)
		}
	}(ctx)
}
