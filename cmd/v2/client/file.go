package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/terrywh/devkit/cmd/v2/stream"
)

type StreamFile struct {
	Path string
}

func (sf *StreamFile) Request() stream.Request {
	return &stream.StreamFileReq{}
}

func (sf *StreamFile) ServeStream(ctx context.Context, r io.Reader, w io.Writer) {
	info, err := os.Stat(sf.Path)
	if err != nil || info.IsDir() {
		log.Println("failed to stat file: ", err)
		return
	}
	file, err := os.Open(sf.Path)
	if err != nil {
		log.Println("failed to open file: ", err)
		return
	}
	defer file.Close()
	log.Println("<client.StreamFile> streaming file: ", sf.Path, info.Size())

	fmt.Fprintf(w, `StreamFile:{"name":"%s","size":%d, "perm":%d}%s`, filepath.Base(sf.Path), info.Size(), info.Mode().Perm(), "\n")
	if size, err := io.Copy(w, file); err != nil || size != info.Size() {
		log.Println("failed to stream file: ", err, "or data corruption")
		return
	}
	if c, ok := w.(io.Closer); ok {
		c.Close() // ShutdownWrite
	}
	rb := bufio.NewReader(r)
	x, err := rb.ReadString('\n')
	if err != nil || x != "done\n" {
		log.Println("failed to stream file: ", x, err)
		return
	}
	log.Println("<client.StreamFile> done.")
}

type FetchFile struct{}

func (ff *FetchFile) ServeStream() {}
