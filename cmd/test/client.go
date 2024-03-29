package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/terrywh/devkit/stream"
	"golang.org/x/term"
)

func RunClient(ctx context.Context, wg *sync.WaitGroup) {
	var c quic.Connection
	var err error
	for {
		c, err = stream.CreateClient()
		if _, ok := err.(*quic.IdleTimeoutError); ok {
			log.Println("retry connecting ...")
			time.Sleep(3 * time.Second)
			continue
		}
		if err != nil {
			log.Println("failed to create client: ", err)
			return
		}
		break
	}
	log.Println("connected ...")
	defer c.CloseWithError(quic.ApplicationErrorCode(0), "done")

	s, err := c.OpenStream()
	if err != nil {
		log.Print(err)
		return
	}
	defer s.Close()

	// StreamFile(s, "./devkit")
	OpenShell(s, "/bin/zsh")
}

func OpenShell(s quic.Stream, shell string) {
	w, h, _ := term.GetSize(int(os.Stdin.Fd()))
	fmt.Fprintf(s, `shell:{"shell":"%s","width":%d,"height":%d}%s`, shell, w, h, "\n")
	log.Println("accepting input ...")
	state, _ := term.MakeRaw(int(os.Stdin.Fd()))

	go io.Copy(s, os.Stdin)
	io.Copy(os.Stdout, s)

	term.Restore(int(os.Stdin.Fd()), state)
	log.Println("close.")
}

func StreamFile(s quic.Stream, path string) {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		log.Println("failed to stat file: ", err)
		return
	}
	file, err := os.Open(path)
	if err != nil {
		log.Println("failed to open file: ", err)
		return
	}
	defer file.Close()
	log.Println("size: ", info.Size())

	fmt.Fprintf(s, `file:{"name":"%s","size":%d, "perm":%d}%s`, filepath.Base(path), info.Size(), info.Mode().Perm(), "\n")
	if size, err := io.Copy(s, file); err != nil || size != info.Size() {
		log.Println("failed to stream file: ", err, " or corruption")
		return
	}
	log.Println("waiting for responses ...")
	s.Close()
	r := bufio.NewReader(s)
	x, err := r.ReadString('\n')
	if err != nil || x != "done\n" {
		log.Println("failed to stream file: ", x, err)
		return
	}
}
