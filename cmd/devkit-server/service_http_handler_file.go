package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/infra"
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
	bash_id, _ := strconv.ParseUint(r.URL.Query().Get("bash_id"), 10, 32)
	shell := DefaultShellHandler.find(int(bash_id))
	if shell == nil {
		handler.Respond(w, entity.ErrSessionNotFound)
		return
	}
	src, err := stream.NewSessionStream(&shell.Server, shell.conn)
	if err != nil {
		handler.Respond(w, fmt.Errorf("failed to create stream: %w", err))
		return
	}
	defer src.CloseRead()
	defer src.CloseWrite()

	sf := entity.StreamFile{}
	if err = json.NewDecoder(r.Body).Decode(&sf); err != nil {
		handler.Respond(w, fmt.Errorf("failed to decode request body: %w", err))
		return
	}
	io.WriteString(src, "/file/pull:")
	if err = app.SendJSON(src, sf); err != nil {
		handler.Respond(w, fmt.Errorf("failed to send json: %w", err))
		return
	}
	rsp := entity.Response{Error: &entity.DefaultErrorCode{}, Data: &sf}
	if err = app.ReadJSON(src.Reader(), &rsp); err != nil {
		handler.Respond(w, fmt.Errorf("failed to read json: %w", err))
		return
	}
	handler.Respond(w, sf)
	if sf.Target.Path != "" { // 指定了目标文件，直接写入文件
		proc := &app.StreamFile{Desc: &sf}
		err = proc.Do(context.Background(), src)
	} else { // 为指定时，在 RESPONSE 流中带回
		sf.Target.Size, err = io.Copy(w, src)
		if err == nil && sf.Target.Size != sf.Source.Size { // 将文件数据透传给 devctl 转写文件
			err = entity.ErrFileCorrupted
		}
	}
	if err != nil {
		infra.Debug("<devkit-server> failed to stream file:", err)
	}
}
