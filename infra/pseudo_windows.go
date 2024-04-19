//go:build windows
// +build windows

package infra

import (
	"context"
	"io"
	"strings"
	"sync"

	"github.com/terrywh/conpty"
)

type WindowsPseudo struct {
	mutex *sync.Mutex
	cpty  *conpty.ConPty
	exit  bool
}

func (wp *WindowsPseudo) Read(recv []byte) (int, error) {
	if wp.exit {
		return 0, io.EOF
	}
	return wp.cpty.Read(recv)
}

func (wp *WindowsPseudo) Write(data []byte) (int, error) {
	return wp.cpty.Write(data)
}

func (wp *WindowsPseudo) Close() (err error) {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()

	wp.exit = true
	if wp.cpty != nil {
		err = wp.cpty.Close()
		wp.cpty = nil
	}
	return
}

func (wp *WindowsPseudo) Resize(cols, rows int) error {
	return wp.cpty.Resize(cols, rows)
}

func (wp *WindowsPseudo) Pid() int {
	return int(wp.cpty.Pid())
}

func StartPty(ctx context.Context, rows, cols int, cmd string, args ...string) (pty Pseudo, err error) {
	wp := &WindowsPseudo{
		mutex: &sync.Mutex{},
	}
	wp.cpty, err = conpty.Start(strings.Join([]string{cmd, strings.Join(args, " ")}, " "),
		conpty.ConPtyDimensions(cols, rows))
	if err != nil {
		return
	}
	go func() {
		wp.cpty.Wait(context.Background()) // 不调用 Wait 时，对 cpty 的 Read 不会停止
		wp.Close()
	}()
	pty = wp
	return
}
