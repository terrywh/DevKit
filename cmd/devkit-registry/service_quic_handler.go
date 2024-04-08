package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/quic-go/quic-go"
	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/stream"
)

type ConnectionTracker interface {
	GetConn(device_id entity.DeviceID) quic.Connection
}

func initServiceQuicHandler(mux *stream.ServeMux, tracker ConnectionTracker) {
	handler := &ServiceQuicHandler{tracker: tracker}
	mux.HandleFunc("/registry/dial", handler.HandleDial)
}

type ServiceQuicHandler struct {
	app.StreamHandlerBase
	tracker ConnectionTracker
}

func (handler *ServiceQuicHandler) HandleDial(ctx context.Context, ss *stream.SessionStream) {
	server := &entity.RemotePeer{}
	json.NewDecoder(ss).Decode(server)
	log.Println("<ServiceQuicHandler.HandleDial> server: ", server.DeviceID, " client: ", ss.RemotePeer().DeviceID)

	client := ss.RemotePeer()
	// P2P 建连：
	// 1. 要求 SERVER 向本测 CLIENT (ss.RemoteAddress()) 发送数据包打洞
	err := handler.dial1RelayToServer(ctx, server, client)
	if err != nil {
		handler.Respond(ss, err)
		return
	}
	// 2. 本地 CLIENT 向远端 SERVER (conn.RemoteAddr()) 建连
	handler.Respond(ss, server)
}

func (handler *ServiceQuicHandler) dial1RelayToServer(ctx context.Context, server *entity.RemotePeer, client *entity.RemotePeer) (err error) {
	conn := handler.tracker.GetConn(server.DeviceID)
	if conn == nil {
		err = entity.ErrSessionNotFound
		return
	}
	server.Address = conn.RemoteAddr().String()

	ss, err := stream.NewSessionStream(server, conn)
	if err != nil {
		return
	}
	err = ss.Invoke(ctx, "/registry/dial", client, nil)
	return
}
