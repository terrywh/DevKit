package handler

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type StreamFile struct {
	Path string `json:"path"`
	Size int64  `json:"size"`
	Perm uint32 `json:"perm"`
}

func (s *StreamFile) ApplyDefaults() {}

func (s *StreamFile) ServeServer(ctx context.Context, r *bufio.Reader, w io.Writer) {
	file, err := os.CreateTemp(filepath.Dir(s.Path), filepath.Base(s.Path))
	log.Println("<StreamFile.ServeServer> writing ", file.Name())
	if err != nil {
		log.Println("<StreamFile.ServeServer> failed to create file: ", err)
		fmt.Fprintln(w, err)
		return
	}
	defer file.Close()
	size, err := io.Copy(file, r)
	if err != nil {
		log.Println("<StreamFile.ServeServer> failed to copy data: ", err)
		fmt.Fprintln(w, err)
		return
	}
	if size != s.Size {
		log.Println("<StreamFile.ServeServer> size mismatched")
		fmt.Fprintln(w, "corruption")
		return
	}
	file.Close()
	os.Chmod(file.Name(), fs.FileMode(s.Perm))
	fmt.Fprintln(w, "done")
}

func (s *StreamFile) ServeClient(ctx context.Context, r *bufio.Reader, w io.Writer) {
	info, err := os.Stat(s.Path)
	if err != nil || info.IsDir() {
		log.Println("<StreamFile.ServeClient> failed to stat file: ", err)
		return
	}
	file, err := os.Open(s.Path)
	if err != nil {
		log.Println("<StreamFile.ServeClient> failed to open file: ", err)
		return
	}
	defer file.Close()
	log.Println("<StreamFile.ServeClient> streaming file: ", s.Path, info.Size())

	fmt.Fprintf(w, `StreamFile:{"name":"%s","size":%d, "perm":%d}%s`, filepath.Base(s.Path), info.Size(), info.Mode().Perm(), "\n")
	if size, err := io.Copy(w, file); err != nil || size != info.Size() {
		log.Println("<StreamFile.ServeClient> failed to stream file: ", err, "or data corruption")
		return
	}
	if c, ok := w.(io.Closer); ok {
		c.Close() // ShutdownWrite
	}
	rb := bufio.NewReader(r)
	x, err := rb.ReadString('\n')
	if err != nil || x != "done\n" {
		log.Println("<StreamFile.ServeClient> failed to stream file: ", x, err)
		return
	}
	log.Println("<StreamFile.ServeClient> done.")
}
