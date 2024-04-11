package main

import (
	"context"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/stream"
)

type QuicFileHandler struct {
	app.StreamHandlerBase
}

func initQuicFileHandler(mux *stream.ServeMux) *QuicFileHandler {
	handler := &QuicFileHandler{}
	mux.HandleFunc("/file/push", handler.HandlePush)
	return handler
}

func (handler *QuicFileHandler) HandlePush(ctx context.Context, src *stream.SessionStream) {
	sf := entity.StreamFilePush{}
	if err := app.ReadJSON(src.Reader(), &sf); err != nil {
		handler.Respond(src, err)
		return
	}

	file, err := os.CreateTemp(filepath.Dir(sf.Path), filepath.Base(sf.Path)+".devkit_tmp_")
	log.Println("<StreamFile.ServeServer> writing ", file.Name())
	if err != nil {
		log.Println("<StreamFile.ServeServer> failed to create file: ", err)
		handler.Respond(src, err)
		return
	}
	defer file.Close()
	size, err := io.Copy(file, src)
	if err != nil {
		log.Println("<StreamFile.ServeServer> failed to copy data: ", err)
		handler.Respond(src, err)
		return
	}
	if size != sf.Size {
		handler.Respond(src, entity.ErrFileCorrupted)
		return
	}
	file.Close()
	os.Chmod(file.Name(), fs.FileMode(sf.Perm))

	if !sf.Override && handler.exists(sf.Path) {
		handler.Respond(src, entity.ErrFileExisted)
		return
	}
	os.Rename(file.Name(), sf.Path)
	handler.Respond(src, nil)
}

func (handler *QuicFileHandler) exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
