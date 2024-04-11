package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
)

type Handler interface {
	InitFlag(fs *flag.FlagSet)
	Do(ctx context.Context) error
}

type HandlerBase struct {
	addr string
}

func (self HandlerBase) HTTPPost(path string, req interface{}) (rsp *http.Response, err error) {
	var payload []byte
	if payload, err = json.Marshal(req); err != nil {
		return
	}
	body := bytes.NewBuffer(payload)
	return http.DefaultClient.Post(
		fmt.Sprintf("http://%s%s", self.addr, path),
		"application/json",
		body,
	)
}
