package k8s

import (
	"context"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/creack/pty"
)

type Session struct {
	Req  Request
	path string
	conf string // *.kubeconfig

	proc *exec.Cmd
	file *os.File // pty
}

// Start
func (s *Session) Start(ctx context.Context) (err error) {
	s.proc = exec.CommandContext(ctx, s.path, "--kubeconfig", s.conf, "exec", "-n", s.Req.Namespace, "-it", s.Req.Pod, "--", "bash")
	s.file, err = pty.StartWithSize(s.proc, &pty.Winsize{Rows: uint16(s.Req.Rows), Cols: uint16(s.Req.Cols)})

	return
}

func (s *Session) Serve(ctx context.Context) (err error) {
	if s.Req.Command != "" {
		time.Sleep(100 * time.Millisecond)
		io.WriteString(s.file, s.Req.Command)
		io.WriteString(s.file, "\r\r")
		time.Sleep(400 * time.Millisecond)
	}
	return
}

func (s *Session) Write(data []byte) (int, error) {
	return s.file.Write(data)
}

func (s *Session) Read(data []byte) (int, error) {
	return s.file.Read(data)
}

func (s *Session) Close() (err error) {
	if s.file == nil {
		return
	}

	io.WriteString(s.file, "\rexit\r")
	time.Sleep(time.Second)
	err = s.file.Close()
	stop := time.AfterFunc(10*time.Second, func() {
		s.proc.Cancel()
	})
	s.proc.Wait()
	s.file = nil
	stop.Stop() // 进程已经停止
	return
}

func (s *Session) Resize(rows, cols int) {
	s.Req.Cols = cols
	s.Req.Rows = rows
	pty.Setsize(s.file, &pty.Winsize{Rows: uint16(rows), Cols: uint16(cols)})
}

func (s *Session) GetSize() (rows, cols int) {
	return s.Req.Rows, s.Req.Cols
}
