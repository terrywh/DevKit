package main

import (
	"fmt"

	"github.com/terrywh/devkit/handler"
	"github.com/terrywh/devkit/stream"
)

func newQuicService() (qs *stream.Server) {
	var err error
	mux := stream.NewServeMux()
	qs, err = stream.DefaultTransport.CreateServer(stream.ServerOptions{
		Handler:   mux,
		Authorize: authorize,
	})
	if err != nil {
		panic(fmt.Sprint("failed to create server: ", err))
	}
	handler.NewServerShellHandler(mux)
	return
}

func authorize(hash string) bool {
	for _, auth := range DefaultConfig.Authorize {
		if hash == auth {
			return true
		}
	}
	return false
}
