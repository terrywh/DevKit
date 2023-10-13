package main

import (
	"bytes"
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

func (s *BashSetup) Serve(ctx context.Context) (err error) {
	log.Println("<bash-setup> preparing ...")
	time.Sleep(2 * time.Second)
	log.Println("<bash-setup> installing ...")
	s.install()
	log.Println("<bash-setup> done.")
	return
}

func (s *BashSetup) isthere() bool {
	io.WriteString(s.server, "ls -l /usr/local/trzsz/trz &> /dev/null; echo 'wemeet-hybrid-bash-setup:' $?\r")
	time.Sleep(200 * time.Millisecond)
	buffer := make([]byte, 1024)
	for !bytes.Contains(buffer, []byte("wemeet-hybrid-bash-setup: ")) {
		s.server.Read(buffer)
	}
	return !bytes.Contains(buffer, []byte("wemeet-hybrid-bash-setup: 2")) // No such file or directory
}

func (s *BashSetup) install() {
	arch := "amd64"
	path := fmt.Sprintf("/Users/terryhaowu/data/htdocs/github.com/terrywh/devkit/bin/trzsz_linux_%s.tar.gz", arch)
	file, _ := os.Open(path)
	defer file.Close()
	io.WriteString(s.server, fmt.Sprintf("base64 -di > /tmp/trzsz_linux_%s.tar.gz\r", arch))
	// time.Sleep(time.Second)
	e := util.NewRfc2045(s.server)
	io.Copy(e, file)
	e.Close()
	// time.Sleep(time.Second)
	io.WriteString(s.server, "\x04\x04") // Ctrl+D x2
	time.Sleep(100 * time.Millisecond)

	io.WriteString(s.server, fmt.Sprintf("tar x -C /tmp -f /tmp/trzsz_linux_%s.tar.gz\r", arch))
	io.WriteString(s.server, fmt.Sprintf("mv /tmp/trzsz_linux_%s /usr/local/trzsz\r", arch))
	io.WriteString(s.server, fmt.Sprintf("rm -rf /tmp/trzsz_linux_%s.tar.gz\r", arch))
	io.WriteString(s.server, "ln -s /usr/local/trzsz/trz /usr/bin/trz\r")
	io.WriteString(s.server, "ln -s /usr/local/trzsz/tsz /usr/bin/tsz\r")
}