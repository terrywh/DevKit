package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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

func (handler *FileHandler) HandlePull(w http.ResponseWriter, r *http.Request) {
	fp := entity.ServerStreamFilePull{}
	if err := json.NewDecoder(r.Body).Decode(&fp); err != nil {
		handler.Respond(w, fmt.Errorf("failed to decode request body: %w", err))
		return
	}
	shell := DefaultShellHandler.find(fp.Pid)
	if shell == nil {
		handler.Respond(w, entity.ErrSessionNotFound)
		return
	}
	fp.DeviceID = shell.DeviceID
	src, err := stream.NewSessionStream(&shell.RemotePeer, shell.conn)
	if err != nil {
		handler.Respond(w, fmt.Errorf("failed to create stream: %w", err))
		return
	}
	log.Println("<FileHandler.HandlePull> streaming file: ", fp.Path, fp.Size, fp.Perm)
	io.WriteString(src, "/file/pull:")
	if err = app.SendJSON(src, fp.StreamFilePull); err != nil {
		handler.Respond(w, fmt.Errorf("failed to send json: %w", err))
		return
	}
	rsp := entity.HttpResponse{Data: &fp.StreamFile}
	if err = app.ReadJSON(src.Reader(), &rsp); err != nil {
		handler.Respond(w, fmt.Errorf("failed to read json: %w", err))
		return
	}
	handler.Respond(w, fp.StreamFile)
	var size int64
	if size, err = io.Copy(w, src); err != nil || size != fp.Size { // 将文件数据透传给 devctl 转写文件
		log.Println("failed to streaming file: ", err, size, "/", fp.Size)
	}
}
