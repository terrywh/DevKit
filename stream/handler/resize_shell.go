package handler

import (
	"bufio"
	"context"
	"io"
)

type ResizeShell struct {
	Cols int `json:"cols"`
	Rows int `json:"rows"`
}

func (o *ResizeShell) ApplyDefaults() {
	if o.Rows < 16 {
		o.Rows = 16
	}
	if o.Cols < 96 {
		o.Cols = 96
	}
}

func (o *ResizeShell) ServeServer(ctx context.Context, r *bufio.Reader, w io.Writer) {
	o.ApplyDefaults()
}

func (o *ResizeShell) ServeClient(ctx context.Context, r *bufio.Reader, w io.Writer) {
	o.ApplyDefaults()
}
