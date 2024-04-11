package app

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/terrywh/devkit/entity"
)

type HttpHandlerBase struct {
	StreamHandlerBase
}

func (handler HttpHandlerBase) Respond(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	handler.StreamHandlerBase.Respond(w, data)
}

type StreamHandlerBase struct{}

func (StreamHandlerBase) Respond(w io.Writer, data interface{}) (reply error) {
	top, ok := data.(error)
	if !ok { // 非错误信息
		reply = json.NewEncoder(w).Encode(entity.HttpResponse{
			Error: entity.ErrSuccess,
			Data:  data,
		})
		return
	}
	err := top
	for err != nil { // 尝试找到 ErrorCode 信息
		if ewc, ok := err.(entity.ErrorWithCode); ok {
			reply = json.NewEncoder(w).Encode(entity.HttpResponse{
				Error: entity.HttpError{
					Code: ewc.ErrorCode(), Info: top.Error(),
				},
			})
			return
		}
		err = errors.Unwrap(err)
	}
	// 未知的错误类型
	reply = json.NewEncoder(w).Encode(entity.HttpResponse{
		Error: entity.HttpError{
			Code: entity.ErrUnknown.Code,
			Info: top.Error(),
		},
	})
	return
}
