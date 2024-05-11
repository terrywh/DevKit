package main

import (
	"context"
	"crypto"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/ncruces/zenity"
	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/infra/log"
	"github.com/terrywh/devkit/stream"
)

type QuicFileHandler struct {
	app.StreamHandlerBase
	mgr stream.SessionManager
}

func initFileHandler(mgr stream.SessionManager, mux *stream.ServeMux) *QuicFileHandler {
	handler := &QuicFileHandler{mgr: mgr}
	mux.HandleFunc("/file/list", handler.HandleList)
	mux.HandleFunc("/file/pull", handler.HandlePull)
	mux.HandleFunc("/file/push", handler.HandlePush)
	return handler
}

func (handler *QuicFileHandler) HandleList(ctx context.Context, src *stream.SessionStream) {
	var err error

	var files []string
	if files, err = zenity.SelectFileMultiple(); err != nil {
		handler.Respond(src, err)
		return
	}
	var sf []entity.SelectFile
	for _, file := range files {
		f := entity.SelectFile{Path: file, Expire: time.Now().Add(300 * time.Second).Unix()}
		f.Auth = f.GenSign(DefaultConfig.Get().cert.PrivateKey.(crypto.Signer))
		sf = append(sf, f)
	}
	handler.Respond(src, sf)
}

// 拉取文件（验证授权，经过 HandleList 选择的文件）
func (handler *QuicFileHandler) HandlePull(ctx context.Context, src *stream.SessionStream) {
	var err error
	rf := entity.SelectFile{}
	if err = app.ReadJSON(src.Reader(), &rf); err != nil {
		handler.Respond(src, err)
		return
	}
	if rf.GenSign(DefaultConfig.Get().cert.PrivateKey.(crypto.Signer)) != rf.Auth {
		handler.Respond(src, entity.ErrUnauthorized)
		return
	}
	sf := entity.StreamFile{
		Source: entity.File{
			Path: rf.Path,
		},
	}
	info, err := os.Stat(sf.Source.Path)
	if err != nil {
		handler.Respond(src, err)
		return
	}
	sf.Source.Perm = uint32(info.Mode().Perm())
	sf.Source.Size = info.Size()

	log.InfoContext(ctx, "<devkit-client> stream file: ", sf.Source.Path)
	if err = handler.Respond(src, sf); err != nil {
		return
	}

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

func (handler *QuicFileHandler) HandlePush(ctx context.Context, src *stream.SessionStream) {
	var err error
	sf := entity.StreamFile{}
	if err = app.ReadJSON(src.Reader(), &sf); err != nil {

		handler.Respond(src, err)
		return
	}
	if sf.Target.Path, err = zenity.SelectFileSave(
		zenity.Filename(filepath.Base(sf.Source.Path)),
		zenity.ConfirmOverwrite()); err != nil {
		handler.Respond(src, err)
		return
	}

	log.InfoContext(ctx, "<devkit-client> stream file: ", sf.Target.Path)
	sf.Options.Override = true
	proc := app.NewStreamFile(&sf, false)
	if err = proc.Do(ctx, src); err != nil {
		handler.Respond(src, err)
		return
	}
	handler.Respond(src, nil)
}
