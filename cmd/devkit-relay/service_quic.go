package main

import (
	"fmt"

	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/stream"
)

func newQuicService() (qs *stream.Server) {
	var err error
	mux := stream.NewServeMux()
	tracker := newDiscoveryTracker()
	qs, err = stream.NewServer(&stream.ServerOptions{
		Handler: &stream.DefaultConnectionHandler{
			Handler: mux,
			Tracker: tracker,
		},
		Authorize: func(device_id entity.DeviceID) bool {
			return true // 任何节点均可连接
		},
		Certificate:         DefaultConfig.Get().Server.Certificate,
		PrivateKey:          DefaultConfig.Get().Server.PrivateKey,
		ApplicationProtocol: "devkit",
	})
	if err != nil {
		panic(fmt.Sprint("failed to create server: ", err))
	}
	initServiceQuicHandler(mux, tracker)
	return qs
}
