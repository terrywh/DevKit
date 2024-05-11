package main

import (
	"context"
	"io"
	"sync"

	"github.com/quic-go/quic-go"
	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/infra"
	"github.com/terrywh/devkit/infra/log"
	"github.com/terrywh/devkit/stream"
)

type ShellHandler struct {
	app.StreamHandlerBase
	start map[entity.ShellID]*ServerShell
	mutex *sync.RWMutex
}

type ServerShell struct {
	entity.ServerShell
	cpid int          `json:"-"`
	cpty infra.Pseudo `json:"-"`
	conn quic.Connection
}

var DefaultShellHandler *ShellHandler

func initShellHandler(mux *stream.ServeMux) *ShellHandler {
	handler := &ShellHandler{
		start: make(map[entity.ShellID]*ServerShell),
		mutex: &sync.RWMutex{},
	}
	mux.HandleFunc("/shell/start", handler.HandleStart)
	mux.HandleFunc("/shell/resize", handler.HandleResize)
	// TODO cleanup
	DefaultShellHandler = handler
	return handler
}

func (h *ShellHandler) put(e *ServerShell) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.start[e.Shell.ID] = e
}

func (h *ShellHandler) del(e *ServerShell) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	delete(h.start, e.Shell.ID)
}

func (h *ShellHandler) get(id entity.ShellID) *ServerShell {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.start[id]
}

func (h *ShellHandler) find(pid int) *ServerShell {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for _, shell := range h.start {
		if shell.cpid == pid {
			return shell
		}
	}
	return nil
}

func (hss *ShellHandler) HandleStart(ctx context.Context, src *stream.SessionStream) {
	var err error
	e := &ServerShell{}
	if err = app.ReadJSON(src.Reader(), &e); err != nil {
		hss.Respond(src, err)
		return
	}
	e.Shell.ApplyDefaults()

	e.cpty, err = infra.StartPty(ctx, e.Shell.Row, e.Shell.Col, e.Shell.Cmd[0], e.Shell.Cmd[1:]...)
	if err != nil {
		log.Warn("<devkit-server> failed to start shell (start): ", err)
		hss.Respond(src, err)
		return
	}
	defer e.cpty.Close()

	log.Trace("<devkit-server> shell started: ", &e.cpty)
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		io.Copy(src, e.cpty)
		src.CloseWrite()
	}()
	go func() {
		defer wg.Done()
		io.Copy(e.cpty, src)
		e.cpty.Close()
		// src.CloseRead()
	}()
	e.conn = src.Conn
	e.cpid = e.cpty.Pid()
	hss.put(e)

	wg.Wait()

	hss.del(e)
	log.Trace("<devkit-server> shell closed: ", &e.cpty)
}

func (hss *ShellHandler) HandleResize(ctx context.Context, src *stream.SessionStream) {
	e1 := &ServerShell{}
	if err := app.ReadJSON(src.Reader(), &e1); err != nil {
		hss.Respond(src, err)
		return
	}

	e2 := hss.get(e1.Shell.ID)
	if e2 == nil {
		hss.Respond(src, entity.ErrSessionNotFound)
		return
	}
	e2.Shell.Col = e1.Shell.Col
	e2.Shell.Row = e1.Shell.Row
	if err := e2.cpty.Resize(e2.Shell.Col, e2.Shell.Row); err != nil {
		hss.Respond(src, err)
		return
	}
	hss.Respond(src, nil)
}
