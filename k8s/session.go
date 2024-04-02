package k8s

import (
	"context"
	"io"
	"time"
)

type Pseudo interface {
	io.ReadWriteCloser
	Resize(cols, rows int) error
}

type Session struct {
	Req  Request
	path string
	conf string // *.kubeconfig

	cpty Pseudo
}

func (s *Session) Start(ctx context.Context) (err error) {
	s.cpty, err = StartSession(ctx, s)
	return
}

func (s *Session) Serve(ctx context.Context) (err error) {
	if err = s.proc.Start(); err != nil {
		return err
	}
	if s.Req.Command != "" {
		time.Sleep(100 * time.Millisecond)
		io.WriteString(s.cpty, s.Req.Command)
		io.WriteString(s.cpty, "\r\r")
		time.Sleep(400 * time.Millisecond)
	}
	return
}

func (s *Session) Write(data []byte) (int, error) {
	return s.cpty.Write(data)
}

func (s *Session) Read(data []byte) (int, error) {
	return s.cpty.Read(data)
}

func (s *Session) Close() (err error) {
	if s.cpty == nil {
		return
	}

	io.WriteString(s.cpty, "\rexit\r")
	time.Sleep(time.Second)
	err = s.cpty.Close()
	s.cpty = nil
	return
}

func (s *Session) Resize(rows, cols int) {
	s.Req.Cols = cols
	s.Req.Rows = rows
	s.cpty.Resize(cols, rows)
}

func (s *Session) GetSize() (rows, cols int) {
	return s.Req.Rows, s.Req.Cols
}
