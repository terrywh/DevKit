package main

import (
	"fmt"

	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/stream"
)

func newQuicService() (qs *stream.Server) {
	var err error
	mux := stream.NewServeMux()
	qs, err = stream.NewServer(&stream.ServerOptions{
		Handler: &stream.DefaultConnectionHandler{
			Handler: mux,
			Tracker: stream.NewDefaultConnectionTracker(),
		},
		Authorize:           onAuthorize,
		Certificate:         DefaultConfig.Get().Server.Certificate,
		PrivateKey:          DefaultConfig.Get().Server.PrivateKey,
		ApplicationProtocol: "devkit",
	})
	if err != nil {
		panic(fmt.Sprint("failed to create server: ", err))
	}
	initShellHandler(mux)
	initDeviceHandler(mux)
	initQuicFileHandler(mux)
	return
}

func onAuthorize(hash entity.DeviceID) bool {
	for _, auth := range DefaultConfig.Get().Server.Authorized {
		if hash == entity.DeviceID(auth) {
			return true
		}
	}
	return false
}
