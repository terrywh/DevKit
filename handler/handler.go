package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/terrywh/devkit/entity"
)

type HttpHandlerBase struct{}

func (HttpHandlerBase) Respond(w http.ResponseWriter, r *entity.HttpResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(r)
}

func (HttpHandlerBase) Success(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entity.HttpResponse{
		Error: entity.HttpError{
			Code: 0, Info: "success",
		},
		Data: data,
	})
}

func (HttpHandlerBase) Failure(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	if ewc, ok := err.(entity.ErrorWithCode); ok {
		json.NewEncoder(w).Encode(entity.HttpResponse{
			Error: entity.HttpError{
				Code: ewc.ErrorCode(), Info: ewc.Error(),
			},
		})
	} else {
		json.NewEncoder(w).Encode(entity.HttpResponse{
			Error: entity.HttpError{
				Code: entity.ErrUnknown.Code,
				Info: err.Error(),
			},
		})
	}
}

type StreamHandlerBase struct{}

func (StreamHandlerBase) Respond(w io.Writer, r *entity.HttpResponse) {
	json.NewEncoder(w).Encode(r)
}

func (StreamHandlerBase) Success(w io.Writer, data interface{}) {
	json.NewEncoder(w).Encode(entity.HttpResponse{
		Error: entity.HttpError{
			Code: 0, Info: "success",
		},
		Data: data,
	})
}

func (StreamHandlerBase) Failure(w io.Writer, err error) {
	if ewc, ok := err.(entity.ErrorWithCode); ok {
		json.NewEncoder(w).Encode(entity.HttpResponse{
			Error: entity.HttpError{
				Code: ewc.ErrorCode(), Info: ewc.Error(),
			},
		})
	} else {
		json.NewEncoder(w).Encode(entity.HttpResponse{
			Error: entity.HttpError{
				Code: entity.ErrUnknown.Code,
				Info: err.Error(),
			},
		})
	}
}
