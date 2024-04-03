package main

import (
	"flag"

	"github.com/terrywh/devkit/stream"
)

type Options struct {
	http string
	addr string
	cert string
	pkey string
}

var DefaultOptions Options

func main() {
	flag.StringVar(&DefaultOptions.http, "http", "127.0.0.1:18080", "serve web ui on this address")
	flag.StringVar(&DefaultOptions.addr, "bind", "0.0.0.0:18080", "serve QUIC stream on this address")
	flag.StringVar(&DefaultOptions.cert, "cert", "./var/cert/server.crt", "certificate to use for QUIC (server only)")
	flag.StringVar(&DefaultOptions.pkey, "pkey", "./var/pkey/server.key", "private key to use for QUIC (server only)")
	flag.Parse()

	stream.InitTransport(stream.TransportOptions{LocalAddress: DefaultOptions.addr})
	defer stream.DefaultTransport.Close()

}
