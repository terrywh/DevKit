package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"sync"

	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/infra"
	"github.com/terrywh/devkit/stream"
)

type ServiceQuicHandlerShell struct {
	app.StreamHandlerBase
	start map[entity.ShellID]*ServerShell
	mutex *sync.RWMutex
}

type ServerShell struct {
	entity.RemoteShell
	cpty infra.Pseudo `json:"-"`
}

func initServiceQuicShellHandler(mux *stream.ServeMux) {
	handler := &ServiceQuicHandlerShell{
		start: make(map[entity.ShellID]*ServerShell),
		mutex: &sync.RWMutex{},
	}
	mux.HandleFunc("/shell/start", handler.HandleStart)
	mux.HandleFunc("/shell/resize", handler.HandleResize)
	// TODO cleanup

}

func (h *ServiceQuicHandlerShell) put(e *ServerShell) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.start[e.ShellId] = e
}

func (h *ServiceQuicHandlerShell) del(e *ServerShell) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	delete(h.start, e.ShellId)
}

func (h *ServiceQuicHandlerShell) get(id entity.ShellID) *ServerShell {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.start[id]
}

func (hss *ServiceQuicHandlerShell) HandleStart(ctx context.Context, ss *stream.SessionStream) {
	log.Println("<ServiceQuicHandlerShell.HandleStart> device =", ss.RemotePeer().DeviceID)
	var err error
	e := &ServerShell{}
	json.NewDecoder(ss).Decode(&e)
	e.ApplyDefaults()

	e.cpty, err = infra.StartPty(ctx, e.Rows, e.Cols, e.ShellCmd[0], e.ShellCmd[1:]...)
	if err != nil {
		log.Println("<HandlerServerShell.HandleStream> failed to open pty shell: ", err)
		hss.Respond(ss, err)
		return
	}
	defer e.cpty.Close()

	hss.put(e)
	defer hss.del(e)

	log.Println("<ServerShellHandler> shell: ", &e.cpty, " started ...")
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		io.Copy(ss, e.cpty)
		ss.CloseRead()
		ss.Close()
	}()
	go func() {
		defer wg.Done()
		io.Copy(e.cpty, ss)
		e.cpty.Close()
	}()
	wg.Wait()
	log.Println("<ServerShellHandler> shell: ", &e.cpty, " closed.")
}

func (hss *ServiceQuicHandlerShell) HandleResize(ctx context.Context, ss *stream.SessionStream) {
	e1 := &ServerShell{}
	json.NewDecoder(ss).Decode(&e1)

	e2 := hss.get(e1.ShellId)
	if e2 == nil {
		hss.Respond(ss, entity.ErrSessionNotFound)
		return
	}
	e2.Cols = e1.Cols
	e2.Rows = e1.Rows
	if err := e2.cpty.Resize(e2.Cols, e2.Rows); err != nil {
		log.Println("<ServerShellHandler.HandleResize> failed to resize cpty: ", err)
		hss.Respond(ss, err)
		return
	}
	hss.Respond(ss, nil)
}
