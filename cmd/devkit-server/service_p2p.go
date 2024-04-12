package main

import (
	"github.com/terrywh/devkit/stream"
)

func newP2PService() (svc *stream.Client) {
	mux := stream.NewServeMux()
	svc = stream.NewClient(&stream.ClientOptions{
		Handler: &stream.DefaultConnectionHandler{
			Tracker: stream.NewDefaultConnectionTracker(),
			Handler: mux,
		},
		DialOptions: stream.DialOptions{
			Address:             DefaultConfig.Get().Relay.Address,
			Certificate:         DefaultConfig.Get().Server.Certificate,
			PrivateKey:          DefaultConfig.Get().Server.PrivateKey,
			ApplicationProtocol: "devkit",
		},
	})
	initServiceP2PHandler(mux)
	return
}
