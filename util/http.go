package util

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
)

func HttpGet(ctx context.Context, url string) (out string, err error) {
	var rsp *http.Response
	if rsp, err = http.Get(url); err != nil {
		return
	}
	defer rsp.Body.Close()
	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, rsp.Body)
	return buf.String(), err
}

func HttpDownload(ctx context.Context, url string, path string) (err error) {
	var rsp *http.Response
	if rsp, err = http.Get(url); err != nil {
		return
	}
	defer rsp.Body.Close()
	var file *os.File
	if file, err = os.Create(path); err != nil {
		return
	}
	defer file.Close()
	_, err = io.Copy(file, rsp.Body)
	return
}
