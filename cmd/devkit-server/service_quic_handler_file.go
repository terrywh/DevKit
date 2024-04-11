package main

import (
	"context"

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
	sf := entity.StreamFile{}
	if err := app.ReadJSON(src.Reader(), &sf); err != nil {
		handler.Respond(src, err)
		return
	}

	proc := app.StreamFile{Desc: &sf}
	handler.Respond(src, proc.Do(ctx, src))
}
