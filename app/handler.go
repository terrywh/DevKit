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
		reply = json.NewEncoder(w).Encode(entity.Response{
			// Error: entity.ErrSuccess,
			Data: data,
		})
		return
	}
	err := top
	for err != nil { // 尝试找到 ErrorCode 信息
		if ewc, ok := err.(entity.ErrorCode); ok {
			reply = json.NewEncoder(w).Encode(entity.Response{
				Error: &entity.DefaultErrorCode{
					Code: ewc.ErrCode(), Info: top.Error(),
				},
			})
			return
		}
		err = errors.Unwrap(err)
	}
	// 未知的错误类型
	reply = json.NewEncoder(w).Encode(entity.Response{
		Error: &entity.DefaultErrorCode{
			Code: entity.ErrUnknown.Code, Info: top.Error(),
		},
	})
	return
}
