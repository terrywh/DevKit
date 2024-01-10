package ssh

import (
	"context"
	"io"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)


type Session struct {
	Req  Request
	session *ssh.Session
	
	cor *io.PipeReader
	cow *io.PipeWriter

	sor io.Reader
	ser io.Reader
	siw io.WriteCloser
}

func (s *Session) Start(ctx context.Context) (err error) {
	if s.sor, err = s.session.StdoutPipe(); err != nil {
		return err
	}
	if s.ser, err = s.session.StderrPipe(); err != nil {
		return err
	}
	if s.siw, err = s.session.StdinPipe(); err != nil {
		return err
	}
	s.cor, s.cow = io.Pipe()
	return
}

func (s *Session) Serve(ctx context.Context) (err error) {
	var wg sync.WaitGroup
	
	wg.Add(1)
	go func() {
		defer wg.Done()
		io.Copy(s.cow, s.sor) // stdout => output
		s.cow.Close()
	} ()
	wg.Add(1)
	go func() {
		defer wg.Done()
		io.Copy(s.cow, s.ser) // stderr => output
	} ()
	ch := make(chan error)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := s.session.Shell(); err != nil {
			ch <- err
			return
		}
		ch <- s.session.Wait()
	} ()

	// go func(ctx context.Context) { // 此协程用于保活，不参与结束流程等待
	// 	timer := time.NewTicker(25 * time.Second)
	// 	defer timer.Stop()
	// SERVING:
	// 	for {
	// 		select {
	// 		case <- ctx.Done():
	// 			break SERVING
	// 		case <- timer.C:
	// 			s.session.SendRequest("keepalive", false, nil)
	// 		}
	// 	}
	// } (ctx)
	
	if s.Req.Command != "" {
		time.Sleep(500 * time.Millisecond)
		io.WriteString(s.siw, s.Req.Command)
		io.WriteString(s.siw, "\r\r")
	}
	
	select {
	case <- ctx.Done():
		s.cow.Close()
		s.session.Close()
		return ctx.Err()
	case err = <- ch:
		return
	}
}

func (s *Session) Read(data []byte) (int, error) {
	return s.cor.Read(data)
}

func (s *Session) Write(data []byte) (int, error) {
	return s.siw.Write(data)
}

func (s *Session) Close() (err error) {
	err = s.siw.Close()
	s.session.Close()
	return
}

func (s *Session) Resize(rows, cols int) {
	s.Req.Rows = rows
	s.Req.Cols = cols
	s.session.WindowChange(rows, cols)
}

func (s *Session) GetSize() (rows, cols int) {
	return s.Req.Rows, s.Req.Cols
}
