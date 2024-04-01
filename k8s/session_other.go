//go:build !windows
// +build !windows

package k8s

import (
	"context"
	"os"
	"os/exec"
	"time"

	"github.com/creack/pty"
)

type UnixPseudo struct {
	proc *exec.Cmd
	file *os.File
}

func (up UnixPseudo) Read(recv []byte) (int, error) {
	return up.file.Read(recv)
}

func (up UnixPseudo) Write(data []byte) (int, error) {
	return up.file.Write(data)
}

func (up UnixPseudo) Close() (err error) {
	err = up.file.Close()
	up.proc.WaitDelay = 3 * time.Second
	up.proc.Wait()
	return
}
func (up UnixPseudo) Resize(cols, rows int) error {
	return pty.Setsize(up.file, &pty.Winsize{Cols: uint16(cols), Rows: uint16(rows)})
}

func StartSession(ctx context.Context, s *Session) (Pseudo, error) {
	var up UnixPseudo
	var err error
	up.proc = exec.CommandContext(ctx, s.path, "--kubeconfig", s.conf, "exec", "-n", s.Req.Namespace, "-it", s.Req.Pod, "--", "bash")
	up.file, err = pty.StartWithSize(up.proc, &pty.Winsize{Rows: uint16(s.Req.Rows), Cols: uint16(s.Req.Cols)})
	if err != nil {
		return nil, err
	}
	return up, nil
}
