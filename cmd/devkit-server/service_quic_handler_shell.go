package main

import (
	"context"
	"io"
	"log"
	"sync"

	"github.com/quic-go/quic-go"
	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/infra"
	"github.com/terrywh/devkit/stream"
)

type ShellHandler struct {
	app.StreamHandlerBase
	start map[entity.ShellID]*ServerShell
	mutex *sync.RWMutex
}

type ServerShell struct {
	entity.RemoteShell
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
	h.start[e.ShellId] = e
}

func (h *ShellHandler) del(e *ServerShell) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	delete(h.start, e.ShellId)
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
	log.Println("<ServiceQuicHandlerShell.HandleStart> device =", src.RemotePeer().DeviceID)
	var err error
	e := &ServerShell{}
	if err = src.Pull(&e); err != nil {
		hss.Respond(src, err)
		return
	}
	e.ApplyDefaults()

	e.cpty, err = infra.StartPty(ctx, e.Rows, e.Cols, e.ShellCmd[0], e.ShellCmd[1:]...)
	if err != nil {
		log.Println("<HandlerServerShell.HandleStream> failed to open pty shell: ", err)
		hss.Respond(src, err)
		return
	}
	defer e.cpty.Close()

	log.Println("<ServerShellHandler> shell: ", &e.cpty, " started ...")
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		io.Copy(src, e.cpty)
		src.CloseRead()
		src.Close()
	}()
	go func() {
		defer wg.Done()
		io.Copy(e.cpty, src)
		e.cpty.Close()
	}()
	e.conn = src.Conn
	e.cpid = e.cpty.Pid()
	hss.put(e)

	wg.Wait()

	hss.del(e)
	log.Println("<ServerShellHandler> shell: ", &e.cpty, " closed.")
}

func (hss *ShellHandler) HandleResize(ctx context.Context, src *stream.SessionStream) {
	e1 := &ServerShell{}
	if err := src.Pull(&e1); err != nil {
		hss.Respond(src, err)
		return
	}

	e2 := hss.get(e1.ShellId)
	if e2 == nil {
		hss.Respond(src, entity.ErrSessionNotFound)
		return
	}
	e2.Cols = e1.Cols
	e2.Rows = e1.Rows
	if err := e2.cpty.Resize(e2.Cols, e2.Rows); err != nil {
		log.Println("<ServerShellHandler.HandleResize> failed to resize cpty: ", err)
		hss.Respond(src, err)
		return
	}
	hss.Respond(src, nil)
}
