//go:build !windows
// +build !windows

package util

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

func StartPty(ctx context.Context, rows, cols int, cmd string, args ...string) (Pseudo, error) {
	var uerr error
	var upty UnixPseudo
	upty.proc = exec.CommandContext(ctx, cmd, args...)
	upty.file, uerr = pty.StartWithSize(upty.proc, &pty.Winsize{Rows: uint16(rows), Cols: uint16(cols)})
	if uerr != nil {
		return nil, uerr
	}
	return upty, nil
}
