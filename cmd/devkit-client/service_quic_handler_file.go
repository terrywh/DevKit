package main

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/ncruces/zenity"
	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/stream"
)

type FileHandler struct {
	app.StreamHandlerBase
	mgr stream.SessionManager
}

func initFileHandler(mgr stream.SessionManager, mux *stream.ServeMux) *FileHandler {
	handler := &FileHandler{mgr: mgr}
	mux.HandleFunc("/file/pull", handler.HandlePull)
	return handler
}

func (handler *FileHandler) HandlePull(ctx context.Context, src *stream.SessionStream) {
	sf := entity.FilePull{}

	if err := app.ReadJSON(src.Reader(), &sf); err != nil {
		handler.Respond(src, err)
		return
	}
	path, err := zenity.SelectFile()
	if err != nil {
		handler.Respond(src, err)
		return
	}

	info, err := os.Stat(path)
	if err != nil {
		handler.Respond(src, err)
		return
	}

	// dst, err := handler.mgr.Acquire(ctx, &src.Peer)
	// if err != nil {
	// 	handler.Respond(src, err)
	// 	return
	// }
	// defer dst.Close()

	file, err := os.Open(path)
	if err != nil {
		handler.Respond(src, err)
		return
	}
	defer file.Close()
	sf.Path = filepath.Join(sf.Path, filepath.Base(path))
	sf.Perm = uint32(info.Mode().Perm())
	sf.Size = info.Size()

	log.Println("<StreamFile.ServeClient> streaming file: ", path, " => ", sf.Path, sf.Size, sf.Perm)

	if err = handler.Respond(src, sf); err != nil {
		log.Println("<FileHandler.HandlePull> failed to respond: ", err)
		return
	}

	if size, err := io.Copy(src, file); err != nil || size != sf.Size {
		log.Println("<StreamFile.ServeClient> failed to stream file: ", err, "or data corruption")
		handler.Respond(src, entity.ErrFileCorrupted)
		return
	}
	src.Close() // 关闭写（发送完毕）
	x := entity.HttpResponse{}
	if err = app.ReadJSON(src.Reader(), &x); err == nil {
		err = x.Error
	}
	if err != nil {
		log.Println("<StreamFile.ServeClient> failed to stream file: ", err)
		handler.Respond(src, err)
		return
	}
	handler.Respond(src, sf)
	log.Println("<StreamFile.ServeClient> done.")
}
