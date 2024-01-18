package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/terrywh/devkit/util"
)

type BashSetup struct {
	server BashBackend
}

func (s *BashSetup) Serve(ctx context.Context, arch string) (err error) {
	log.Println("<bash-setup> preparing ...")
	time.Sleep(2 * time.Second)
	log.Println("<bash-setup> installing ...")
	// name := "trzsz_1.1.7_linux_x86_64"
	// name := "trzsz_1.1.7_linux_aarch64"
	s.install(fmt.Sprintf("trzsz_1.1.7_linux_%s", arch))
	log.Println("<bash-setup> done.")
	return
}

func (s *BashSetup) install(name string) {
	

	path := fmt.Sprintf("/Users/terryhaowu/data/htdocs/github.com/terrywh/devkit/var/%s.tar.gz", name)
	file, _ := os.Open(path)
	defer file.Close()
	io.WriteString(s.server, fmt.Sprintf("base64 -di > /tmp/%s.tar.gz\r", name))
	// time.Sleep(time.Second)
	e := util.NewRfc2045(s.server)
	io.Copy(e, file)
	e.Close()
	// time.Sleep(time.Second)
	io.WriteString(s.server, "\x04\x04") // Ctrl+D x2
	time.Sleep(100 * time.Millisecond)

	io.WriteString(s.server, fmt.Sprintf("tar x -C /tmp -f /tmp/%s.tar.gz\r", name))
	io.WriteString(s.server, fmt.Sprintf("mv /tmp/%s /usr/local/trzsz\r", name))
	io.WriteString(s.server, fmt.Sprintf("rm -rf /tmp/%s.tar.gz\r", name))
	io.WriteString(s.server, "ln -s /usr/local/trzsz/trz /usr/bin/trz\r")
	io.WriteString(s.server, "ln -s /usr/local/trzsz/tsz /usr/bin/tsz\r")
}