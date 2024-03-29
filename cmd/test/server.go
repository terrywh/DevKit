package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/creack/pty"
	"github.com/terrywh/devkit/stream"
)

type ShellConfig struct {
	Shell  string `json:"shell"`
	Width  uint16 `json:"width"`
	Height uint16 `json:"height"`
}

type SfileConfig struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Perm uint16 `json:"perm"`
}

func RunServer(ctx context.Context, wg *sync.WaitGroup) {
	time.Sleep(10 * time.Second)

	l, err := stream.CreateServer()
	if err != nil {
		log.Print(err)
		return
	}
	defer l.Close()

	log.Println("wait for connections ...")
	for {
		conn, err := l.Accept(ctx)
		if err != nil {
			log.Print(err)
			return
		}
		go func() {
			for {
				s, err := conn.AcceptStream(ctx)
				if err != nil {
					log.Print(err)
					return
				}
				go func() {
					defer s.Close()
					r := bufio.NewReader(s)
					x, err := r.ReadString('\n')
					if err != nil {
						log.Print(err)
						return
					}
					if strings.HasPrefix(x, "shell:") {
						var sc ShellConfig
						json.Unmarshal([]byte(x[6:]), &sc)
						log.Println("starting shell ", sc)
						proc := exec.CommandContext(ctx, sc.Shell)
						file, err := pty.StartWithSize(proc, &pty.Winsize{Rows: uint16(sc.Height), Cols: uint16(sc.Width)})
						if err != nil {
							log.Println("failed to open pty shell: ", err)
							return
						}
						go io.Copy(s, file)
						io.Copy(file, r)
					} else if strings.HasPrefix(x, "file:") {
						var sc SfileConfig
						json.Unmarshal([]byte(x[5:len(x)-1]), &sc)
						file, err := os.CreateTemp("", sc.Name)
						if err != nil {
							log.Println("failed to create file: ", err)
							fmt.Fprintln(s, err)
							return
						}
						log.Println("prepare to transfer file: ", sc, file.Name())
						defer file.Close()
						size, err := io.Copy(file, r)
						if err != nil {
							log.Println("failed to copy data: ", err)
							fmt.Fprintln(s, err)
							return
						}
						if size != sc.Size {
							log.Println("size mismatched")
							fmt.Fprintln(s, "transfer corruption")
							return
						}
						file.Close()
						os.Chmod(file.Name(), fs.FileMode(sc.Perm))
						fmt.Fprintln(s, "done")
					}
				}()
			}
		}()
	}
}
