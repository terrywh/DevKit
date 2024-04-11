package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
)

type Handler interface {
	InitFlag(fs *flag.FlagSet)
	Do(ctx context.Context) error
	Close() error
}

type HandlerBase struct {
	addr string
}

func (handler HandlerBase) HTTPPost(path string, req interface{}) (rsp *http.Response, err error) {
	var payload []byte
	if payload, err = json.Marshal(req); err != nil {
		return
	}
	body := bytes.NewBuffer(payload)
	return http.DefaultClient.Post(
		fmt.Sprintf("http://%s%s", handler.addr, path),
		"application/json",
		body,
	)
}

type HandlerService struct {
	name    string
	handler Handler
}

func (svc HandlerService) Serve(ctx context.Context) {
	err := svc.handler.Do(ctx)
	if err != nil {
		fmt.Printf("error: failed to %s, due to: %s\n", svc.name, err.Error())
		os.Exit(-1)
		return
	}
}

func (svc HandlerService) Close() error {
	return svc.handler.Close()
}
