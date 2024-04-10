package main

import (
	"context"

	"github.com/terrywh/devkit/stream"
)

type FileHandler struct{}

func initFileHandler(mux *stream.ServeMux) *FileHandler {
	handler := &FileHandler{}
	mux.HandleFunc("/file/pull", handler.HandlePull)
	return handler
}

func (handler *FileHandler) HandlePull(ctx context.Context, ss *stream.SessionStream) {
	
}
