package main

import (
	"context"

	"github.com/terrywh/devkit/stream"
)

type QuicService struct {
	mux *stream.ServeMux
}

func newQuicService() (qs *QuicService) {
	qs = &QuicService{mux: stream.NewServeMux()}

	return qs
}

func (qs *QuicService) Serve(ctx context.Context) {

}

func (s *QuicService) Close() error {
	return nil
}
