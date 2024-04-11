package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/stream"
)

type HttpFileHandler struct {
	app.HttpHandlerBase
	mgr stream.SessionManager
}

func initHttpFileHandler(mgr stream.SessionManager, mux *http.ServeMux) *HttpFileHandler {
	handler := &HttpFileHandler{mgr: mgr}
	mux.HandleFunc("/file/push", handler.HandlePush)
	return handler
}

func (handler *HttpFileHandler) HandlePush(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()

	sf := entity.StreamFile{}
	sf.Source.Path = query.Get("source")
	sf.Source.Size = req.ContentLength
	perm, _ := strconv.ParseUint(query.Get("perm"), 8, 32)
	sf.Source.Perm = uint32(perm)
	sf.Target.Path = query.Get("target")
	// push.Target.Size = push.Source.Size
	// push.Target.Perm = push.Source.Perm

	target := entity.Server{}
	target.Address = query.Get("address")
	target.DeviceID = entity.DeviceID(query.Get("device_id"))

	if sf.Source.Path == "" || sf.Target.Path == "" || (target.Address == "" && target.DeviceID == "") {
		handler.Respond(w, entity.ErrInvalidArguments)
		return
	}
	dst, err := handler.mgr.Acquire(context.TODO(), &target)
	if err != nil {
		handler.Respond(w, err)
		return
	}
	if handler.isLocalFile(&sf.Source) {
		info, err := os.Stat(sf.Source.Path)
		if err != nil {
			handler.Respond(w, err)
			return
		}
		sf.Source.Size = info.Size()
		sf.Source.Perm = uint32(info.Mode().Perm())
	}
	// 发送请求
	io.WriteString(dst, "/file/push:")
	if err = app.SendJSON(dst, sf); err != nil {
		handler.Respond(w, err)
		return
	}
	// 传输文件
	var size int64
	if handler.isLocalFile(&sf.Source) {
		file, err := os.Open(sf.Source.Path)
		if err != nil {
			handler.Respond(w, err)
			return
		}
		defer file.Close()
		size, err = io.Copy(dst, file)
	} else { // 直接将请求内容写入目标文件
		size, err = io.Copy(dst, req.Body)
	}
	// 检查文件
	if err != nil || size != sf.Source.Size {
		handler.Respond(w, entity.ErrFileCorrupted)
		return
	}
	handler.Respond(w, nil)
}

func (handler *HttpFileHandler) isLocalFile(file *entity.File) bool {
	return file.Size == 0 // 直接在 REQ 流中传递文件应包含 Content-Length 大小
}
