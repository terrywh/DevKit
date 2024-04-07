package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/stream"
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

type StreamHandlerInvoker struct{}

func (shc StreamHandlerInvoker) Invoke(ctx context.Context, device_id entity.DeviceID, path string, req interface{}) (rsp entity.HttpResponse, err error) {
	s, err := stream.DefaultSessionManager.AcquireStream(ctx, device_id)
	if err != nil {
		log.Println("<ServiceHttpShell.HandleSocket> failed to acquire stream: ", err)
		return
	}
	fmt.Fprintf(s, "%s:", path)
	if err = json.NewEncoder(s).Encode(req); err != nil {
		return
	}
	err = json.NewDecoder(s).Decode(&rsp)
	return
}
