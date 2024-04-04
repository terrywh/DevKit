package handler

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

type ExecuteShell struct {
	Shell   []string `json:"shell"`
	Command string   `json:"command"`

	SrcReader io.Reader `json:"-"`
	DstWriter io.Writer `json:"-"`
}

func (e *ExecuteShell) ApplyDefaults() {
	if e.SrcReader == nil {
		e.SrcReader = os.Stdin
	}
	if e.DstWriter == nil {
		e.DstWriter = os.Stdout
	}
}

func (s *ExecuteShell) ServeServer(ctx context.Context, r *bufio.Reader, w io.Writer) {
	start := time.Now()
	cmd := exec.CommandContext(ctx, s.Shell[0], s.Shell[1:]...)
	cmd.Stdout = w
	cmd.Stderr = w
	err := cmd.Run()
	if ee, ok := err.(*exec.ExitError); ok {
		fmt.Fprintf(w, "\n{\"data\":{\"exit\":%d,\"cost\":%.2fs}}\n", ee.ExitCode(), time.Since(start).Seconds())
	} else if err != nil {
		fmt.Fprintln(w, "\n{\"error\":{\"code\":%d,\"info\":\"%s\"}}\n", err)
	} else {
		fmt.Fprintf(w, "\n{\"exit\":0,\"cost\":%.2fs}\n", time.Since(start).Seconds())
	}
}

func (s *ExecuteShell) ServeClient(ctx context.Context, r *bufio.Reader, w io.Writer) {
	io.WriteString(w, "ExecuteShell:")
	e := json.NewEncoder(w)
	e.Encode(s)

	io.Copy(s.DstWriter, r)
}
