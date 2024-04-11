package app

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/stream"
)

func ReadJSON(src *bufio.Reader, dst interface{}) (err error) {
	var payload []byte
	if payload, err = src.ReadBytes(byte('\n')); err != nil {
		return
	}
	err = json.Unmarshal(payload, dst)
	return
}

func SendJSON(dst io.Writer, src interface{}) (err error) {
	return json.NewEncoder(dst).Encode(src)
}

func Invoke(ctx context.Context, src *stream.SessionStream, path string,
	req interface{}, rsp interface{}) (err error) {
	defer src.CloseWrite()

	if _, err = fmt.Fprintf(src, "%s:", path); err != nil {
		return
	}
	if err = SendJSON(src, req); err != nil {
		return
	}
	r := entity.Response{Error: &entity.DefaultErrorCode{}, Data: rsp}
	// Decoder 可能读取了更后面的内容，可以使用（响应内容已结束）
	// err = json.NewDecoder(ss).Decode(&r)
	if err = ReadJSON(src.Reader(), &r); err != nil {
		return
	}

	if r.Error.Code > 0 {
		err = r.Error
		return
	}
	return nil
}
