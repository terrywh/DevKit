package client

import (
	"context"
	"flag"
	"log"
	"net"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
)

type Options struct {
	Addr string
	File string
}

type Client struct {
	handler map[string]stream.ClientHandler
	options Options
}

func New() *Client {
	cli := &Client{
		handler: make(map[string]stream.ClientHandler),
	}
	cli.handler["StreamFile"] = &StreamFile{}
	return cli
}

func (cli *Client) ParseFlags() {
	flag.StringVar(&cli.options.File, "file", "", "file path to transfer to server")
	flag.StringVar(&cli.options.Addr, "addr", "127.0.0.1:12345", "connect to server, eg. 127.0.0.1:12345")
}

// func (cli *Client) Handle(name string, handle Handler) {
// 	cli.handler[name] = handle
// }

func (cli *Client) Serve(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	var c quic.Connection
	var err error
	for {
		addr, _ := net.ResolveUDPAddr("udp", cli.options.Addr)
		c, err = stream.Connect(ctx, addr)
		if _, ok := err.(*quic.IdleTimeoutError); ok {
			log.Println("<client.Client> retry connect in 3s ...")
			time.Sleep(3 * time.Second)
			continue
		}
		if err != nil {
			log.Println("<client.Client> failed to connect: ", err)
			return
		}
		break
	}
	log.Println("<client.Client> connected ...")
	defer c.CloseWithError(quic.ApplicationErrorCode(0), "done")

	s, err := c.OpenStream()
	if err != nil {
		log.Print(err)
		return
	}
	defer s.Close()

	// sf := &StreamFile{Path: cli.options.File}
	// sf.ServeStream(context.Background(), s, s)
	OpenShell(s, "bash")

}
