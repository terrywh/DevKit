package main

import (
	"context"
	"io"
	"log"
	"sync"

	trzsz "github.com/trzsz/trzsz-go/trzsz"
)

type BashBackend interface {
	io.ReadWriteCloser
	Start(ctx context.Context) error
	Resize(rows, cols int)
	GetSize() (rows, cols int)
}

type BashServe interface {
	Serve(ctx context.Context) error
}

type BashShell struct {
	source io.ReadWriteCloser
	server BashBackend
	// filter *trzsz.TrzszFilter
}

type ResizeRequest interface {
	Rows() int
	Cols() int
}

type ResizeRequestImpl struct {
	R int `json:"rows"`
	C int `json:"cols"`
}

func (r ResizeRequestImpl) Rows() int {
	return r.R
}

func (r ResizeRequestImpl) Cols() int {
	return r.C
}

func (s *BashShell) Resize(rows, cols int) {
	// s.filter.SetTerminalColumns(int32(cols))
	s.server.Resize(rows, cols)
}

func (s *BashShell) Serve(ctx context.Context) error {
	log.Println("<shell> starting ...")
	wg := &sync.WaitGroup{}
	defer s.source.Close()
	// 服务端准备工作
	if err := s.server.Start(ctx); err != nil {
		log.Println("failed to start (server): ", err)
	}
	defer s.server.Close()
	// !!! 使用 JS 在浏览器进行对应支持，不再使用服务端 Filter 机制
	// 在 Shell 中支持 trzsz 上传下载
	rsource, wsource := io.Pipe()
	rserver, wserver := io.Pipe()
	_, cols := s.server.GetSize()
	/*s.filter = */ trzsz.NewTrzszFilter(rsource, wserver, s.server, s.server, trzsz.TrzszOptions{
		TerminalColumns: int32(cols),
	})

	wg.Add(1)
	go func() { // 用户输入 (WebSocket) -> Pipe -> Filter -> Server
		defer wg.Done()
		io.Copy(wsource, s.source)
		wsource.Close()
		// io.Copy(s.source, s.server)
		// s.source.Close()
	}()
	wg.Add(1)
	go func() { // 服务输出 Server -> Filter -> Pipe -> (WebSocket)
		defer wg.Done()
		io.Copy(s.source, rserver)
		s.source.Close()
		// io.Copy(s.server, s.source)
		// s.server.Close()
	}()

	if svc, ok := s.server.(BashServe); ok { // 服务端运行 Shell 工作
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := svc.Serve(ctx); err != nil {
				log.Println("failed to serve (server): ", err)
			}
		}()
	}

	wg.Wait()
	log.Println("<shell> closed.")
	return nil
}
