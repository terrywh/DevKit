package app

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/terrywh/devkit/entity"
)

type HttpHandlerBase struct{}

func (HttpHandlerBase) Respond(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if ewc, ok := data.(entity.ErrorWithCode); ok {
		json.NewEncoder(w).Encode(entity.HttpResponse{
			Error: entity.HttpError{
				Code: ewc.ErrorCode(), Info: ewc.Error(),
			},
		})
	} else if err, ok := data.(error); ok && err != nil {
		json.NewEncoder(w).Encode(entity.HttpResponse{
			Error: entity.HttpError{
				Code: entity.ErrUnknown.Code,
				Info: err.Error(),
			},
		})
	} else {
		json.NewEncoder(w).Encode(entity.HttpResponse{
			Error: entity.ErrSuccess,
			Data:  data,
		})
	}
}

type StreamHandlerBase struct{}

func (StreamHandlerBase) Respond(w io.Writer, data interface{}) {
	if ewc, ok := data.(entity.ErrorWithCode); ok {
		json.NewEncoder(w).Encode(entity.HttpResponse{
			Error: entity.HttpError{
				Code: ewc.ErrorCode(), Info: ewc.Error(),
			},
		})
	} else if err, ok := data.(error); ok && err != nil {
		json.NewEncoder(w).Encode(entity.HttpResponse{
			Error: entity.HttpError{
				Code: entity.ErrUnknown.Code,
				Info: err.Error(),
			},
		})
	} else {
		json.NewEncoder(w).Encode(entity.HttpResponse{
			Error: entity.ErrSuccess,
			Data:  data,
		})
	}
}
