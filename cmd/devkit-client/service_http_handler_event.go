package main

import (
	"net/http"

	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/stream"
)

type EventHandler struct {
	app.HttpHandlerBase
}

func initEventHandler(hmux *http.ServeMux, smux *stream.ServeMux) *EventHandler {
	handler := &EventHandler{}
	hmux.HandleFunc("/event/subscribe", handler.HandleSubscribe)
	hmux.HandleFunc("/event/publish", handler.HandlePublish)

	return handler
}

func (handler *EventHandler) HandleSubscribe(rsp http.ResponseWriter, req *http.Request) {
}

func (handler *EventHandler) HandlePublish(rsp http.ResponseWriter, req *http.Request) {

}
