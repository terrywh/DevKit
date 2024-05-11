package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	ps "github.com/mitchellh/go-ps"
	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/infra/log"
)

type Handler interface {
	InitFlag(flagCommand, flagGlobal *flag.FlagSet)
	Do(ctx context.Context) error
	Close() error
}

type HandlerBase struct{}

func (handler HandlerBase) Post(path string, req interface{}) (rsp *http.Response, err error) {
	var body io.Reader
	if r, ok := req.(io.Reader); ok {
		body = r
	} else {
		var payload []byte
		if payload, err = json.Marshal(req); err != nil {
			return
		}
		log.Debug("post body: ", string(payload))
		body = bytes.NewBuffer(payload)
	}
	return http.DefaultClient.Post(
		fmt.Sprintf("http://%s%s", DefaultConfig.Get().Server.Address, path),
		"application/json",
		body,
	)
}

func GetBashPid() (pid int, err error) {
	handler := HandlerBase{}
	// 1. 确认 SERVER 进程
	var rsp *http.Response
	rsp, err = handler.Post("/server/query", nil)
	if err != nil {
		return
	}
	server := entity.Server{}
	if err = app.Read(bufio.NewReader(rsp.Body), &server); err != nil {
		return
	}
	log.Trace("<devkit> GetBashPid() devkit-server:", server)
	// 2. 找到 SERVER 的子、调用方的父进程
	var proc ps.Process
	ppid := os.Getppid()
	for {
		proc, err = ps.FindProcess(ppid)
		if err != nil {
			pid = 0
			return
		}
		log.Trace("<devkit> GetBashPid() find:", proc.Executable(), proc.Pid(), proc.PPid())
		if proc.PPid() < 10 {
			break
		}
		if proc.PPid() == server.Pid {
			pid = proc.Pid()
			return
		}
		ppid = proc.PPid()
	}
	err = entity.ErrSessionNotFound
	return
}

type HandlerService struct {
	name    string
	handler Handler
}

func (svc HandlerService) Serve(ctx context.Context) {
	err := svc.handler.Do(ctx)
	if err != nil {
		fmt.Printf("命令 (%s) 失败: %s\n", svc.name, err.Error())
		os.Exit(-1)
		return
	}
}

func (svc HandlerService) Close() error {
	return svc.handler.Close()
}
