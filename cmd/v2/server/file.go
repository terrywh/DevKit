package server

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"

	"github.com/terrywh/devkit/cmd/v2/stream"
)

type StreamFile struct{}

func (sf StreamFile) ServeStream(ctx context.Context, req stream.Request, r *bufio.Reader, w io.Writer) {
	request := req.(*stream.StreamFileReq)
	if request.Name == "" || request.Size < 1 || request.Perm < 1 {
		fmt.Fprintln(w, "invalid options")
		return
	}

	file, err := os.CreateTemp("", request.Name)
	log.Println("<server.StreamFile> writing ", file.Name())
	if err != nil {
		log.Println("failed to create file: ", err)
		fmt.Fprintln(w, err)
		return
	}
	defer file.Close()
	size, err := io.Copy(file, r)
	if err != nil {
		log.Println("failed to copy data: ", err)
		fmt.Fprintln(w, err)
		return
	}
	if size != request.Size {
		log.Println("size mismatched")
		fmt.Fprintln(w, "transfer corruption")
		return
	}
	file.Close()
	os.Chmod(file.Name(), fs.FileMode(request.Perm))
	fmt.Fprintln(w, "done")
}
