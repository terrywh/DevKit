package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/stream"
)

type FileHandler struct {
	app.HttpHandlerBase
}

func initFileHandler(mux *http.ServeMux) *FileHandler {
	handler := &FileHandler{}
	mux.HandleFunc("/file/pull", handler.HandlePull)
	return handler
}

type ServerStreamFilePull struct {
	entity.StreamFilePull
	Pid int `json:"pid"`
}

func (handler *FileHandler) HandlePull(rsp http.ResponseWriter, req *http.Request) {
	fp := ServerStreamFilePull{}
	if err := json.NewDecoder(req.Body).Decode(&fp); err != nil {
		handler.Respond(rsp, err)
		return
	}
	shell := DefaultShellHandler.find(fp.Pid)
	if shell == nil {
		handler.Respond(rsp, entity.ErrSessionNotFound)
		return
	}
	fp.DeviceID = shell.DeviceID
	ss, err := stream.NewSessionStream(&shell.RemotePeer, shell.conn)
	if err != nil {
		handler.Respond(rsp, err)
	}
	if err := ss.Invoke(context.TODO(), "/file/pull", fp.StreamFilePull, &fp.StreamFilePull); err != nil {
		handler.Respond(rsp, err)
		return
	}
	handler.Respond(rsp, fp.StreamFilePull)
}
