package main

import (
	"context"
	"io"
	"os"

	"github.com/ncruces/zenity"
	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/infra"
	"github.com/terrywh/devkit/stream"
)

type QuicFileHandler struct {
	app.StreamHandlerBase
	mgr stream.SessionManager
}

func initFileHandler(mgr stream.SessionManager, mux *stream.ServeMux) *QuicFileHandler {
	handler := &QuicFileHandler{mgr: mgr}
	mux.HandleFunc("/file/pull", handler.HandlePull)
	return handler
}

func (handler *QuicFileHandler) HandlePull(ctx context.Context, src *stream.SessionStream) {
	var err error
	sf := entity.StreamFile{}
	if err = app.ReadJSON(src.Reader(), &sf); err != nil {
		handler.Respond(src, err)
		return
	}
	if sf.Source.Path, err = zenity.SelectFile(); err != nil {
		handler.Respond(src, err)
		return
	}

	info, err := os.Stat(sf.Source.Path)
	if err != nil {
		handler.Respond(src, err)
		return
	}
	sf.Source.Perm = uint32(info.Mode().Perm())
	sf.Source.Size = info.Size()

	infra.Debug("<devkit-client> streaming file: ", sf.Source.Path)
	if err = handler.Respond(src, sf); err != nil {
		return
	}

	// dst, err := handler.mgr.Acquire(ctx, &src.Peer)
	// if err != nil {
	// 	handler.Respond(src, err)
	// 	return
	// }
	// defer dst.Close()

	file, err := os.Open(sf.Source.Path)
	if err != nil {
		handler.Respond(src, err)
		return
	}
	defer file.Close()

	if size, err := io.Copy(src, file); err != nil || size != sf.Source.Size {
		handler.Respond(src, entity.ErrFileCorrupted)
		return
	}
	src.CloseWrite() // 关闭写（发送完毕）
	handler.Respond(src, sf)
}
