package server

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log"
	"os/exec"

	"github.com/creack/pty"
)

type ShellOptions struct {
	Shell  string `json:"shell"`
	Width  uint16 `json:"width"`
	Height uint16 `json:"height"`
}

func (options *ShellOptions) MakeDefaults() {
	if options.Shell == "" {
		options.Shell = "bash"
	}
	if options.Height < 10 {
		options.Height = 10
	}
	if options.Width < 40 {
		options.Width = 40
	}
}

type OpenShell struct{}

func (os *OpenShell) ServeStream(ctx context.Context, opts string, r *bufio.Reader, w io.Writer) {
	var options ShellOptions
	json.Unmarshal([]byte(opts), &options)
	options.MakeDefaults()

	log.Println("open shell: ", options)
	proc := exec.CommandContext(ctx, options.Shell)
	file, err := pty.StartWithSize(proc, &pty.Winsize{
		Rows: uint16(options.Height),
		Cols: uint16(options.Width),
	})
	if err != nil {
		log.Println("failed to open pty shell: ", err)
		return
	}
	go io.Copy(w, file)
	io.Copy(file, r)
}
