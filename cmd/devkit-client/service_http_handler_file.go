package main

import (
	"bufio"
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

func (handler *HttpFileHandler) HandlePush(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	sf := entity.StreamFile{}
	sf.Options.Override, _ = strconv.ParseBool(query.Get("override"))
	sf.Source.Path = query.Get("source")
	sf.Source.Size = r.ContentLength
	perm, _ := strconv.ParseUint(query.Get("perm"), 8, 32)
	sf.Source.Perm = uint32(perm)
	if sf.Source.Perm == 0 {
		sf.Source.Perm = 0o644
	}
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
	defer dst.CloseWrite()
	rbody := bufio.NewReader(r.Body)
	var fromfile bool
	if _, err = rbody.Peek(1); err != nil {
		fromfile = true
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
	if fromfile {
		var file *os.File
		file, err = os.Open(sf.Source.Path)
		if err != nil {
			handler.Respond(w, err)
			return
		}
		defer file.Close()
		size, err = io.Copy(dst, file)
	} else {
		size, err = io.Copy(dst, rbody) // 尝试直接将请求内容写入目标文件
	}
	dst.CloseWrite() // 文件发送完毕
	// 检查文件
	if err != nil {
		handler.Respond(w, err)
		return
	}
	if sf.Options.Override && size != sf.Source.Size {
		handler.Respond(w, entity.ErrFileCorrupted)
		return
	}
	err = app.Read(dst.Reader(), nil)
	handler.Respond(w, err)
}
