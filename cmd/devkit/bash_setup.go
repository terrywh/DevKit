package main

import (
	"bytes"
	"context"
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
	path := "/Users/terryhaowu/data/htdocs/github.com/terrywh/devkit/bin/trzsz_linux_amd64.tar.gz"
	file, _ := os.Open(path)
	defer file.Close()
	io.WriteString(s.server, "base64 -di > /tmp/trzsz_linux_amd64.tar.gz\r")
	// time.Sleep(time.Second)
	e := util.NewRfc2045(s.server)
	io.Copy(e, file)
	e.Close()
	// time.Sleep(time.Second)
	io.WriteString(s.server, "\x04\x04") // Ctrl+D x2
	time.Sleep(100 * time.Millisecond)

	io.WriteString(s.server, "tar x -C /tmp -f /tmp/trzsz_linux_amd64.tar.gz\r")
	io.WriteString(s.server, "mv /tmp/trzsz_1.1.5_linux_x86_64 /usr/local/trzsz\r")
	io.WriteString(s.server, "rm -rf /tmp/trzsz_linux_amd64.tar.gz\r")
	io.WriteString(s.server, "ln -s /usr/local/trzsz/trz /usr/bin/trz\r")
	io.WriteString(s.server, "ln -s /usr/local/trzsz/tsz /usr/bin/tsz\r")
}