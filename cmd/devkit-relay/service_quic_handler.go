package main

import (
	"context"

	"github.com/quic-go/quic-go"
	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/infra/log"
	"github.com/terrywh/devkit/stream"
)

type ConnectionTracker interface {
	GetConn(device_id entity.DeviceID) quic.Connection
}

func initServiceQuicHandler(mux *stream.ServeMux, tracker ConnectionTracker) {
	handler := &ServiceQuicHandler{tracker: tracker}
	mux.HandleFunc("/relay/dial", handler.HandleDial)
}

type ServiceQuicHandler struct {
	app.StreamHandlerBase
	tracker ConnectionTracker
}

func (handler *ServiceQuicHandler) HandleDial(ctx context.Context, src *stream.SessionStream) {
	server := entity.Server{}
	if err := app.ReadJSON(src.Reader(), &server); err != nil {
		handler.Respond(src, err)
		return
	}
	client := src.RemotePeer()
	log.Info("<devkit-relay> dail ", client.DeviceID, "(", client.Address, ") => ", server.DeviceID, "(", server.Address, ")")

	// P2P 建连：
	// 1. 要求 SERVER 向本测 CLIENT (ss.RemoteAddress()) 发送数据包打洞
	err := handler.dial1RelayToServer(ctx, &server, client)
	if err != nil {
		handler.Respond(src, err)
		return
	}
	// 2. 本地 CLIENT 向远端 SERVER (conn.RemoteAddr()) 建连
	handler.Respond(src, server)
}

func (handler *ServiceQuicHandler) dial1RelayToServer(ctx context.Context,
	server *entity.Server, client *entity.Server) (err error) {
	conn := handler.tracker.GetConn(server.DeviceID)
	if conn == nil {
		err = entity.ErrSessionNotFound
		return
	}
	server.Address = conn.RemoteAddr().String()

	src, err := stream.NewSessionStream(server, conn)
	if err != nil {
		return
	}
	err = app.Invoke(ctx, src, "/relay/dial", client, nil)
	return
}
