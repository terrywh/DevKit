package main

import (
	"flag"

	"github.com/terrywh/devkit/application"
	"github.com/terrywh/devkit/stream"
)

type Options struct {
	http string
	addr string
	cert string
	pkey string
}

var DefaultOption Options

type Config struct {
	Authorize []string `yaml:"authorize"`
}

var DefaultConfig Config

func main() {
	application.InitConfigWatcher("devkit", &DefaultConfig)
	defer application.DefaultConfigWatcher.Close()

	flag.StringVar(&DefaultOption.http, "http", "127.0.0.1:18080", "serve web ui on this address")
	flag.StringVar(&DefaultOption.addr, "bind", "0.0.0.0:18081", "serve QUIC stream on this address")
	flag.StringVar(&DefaultOption.cert, "cert", "./var/cert/server.crt", "certificate to use for QUIC (server only)")
	flag.StringVar(&DefaultOption.pkey, "pkey", "./var/cert/server.key", "private key to use for QUIC (server only)")
	flag.Parse()

	stream.InitTransport(stream.TransportOptions{LocalAddress: DefaultOption.addr})
	defer stream.DefaultTransport.Close()

	sc := application.NewServiceController()
	sc.Start(newQuicService())
	sc.WaitForSignal()
	sc.Shutdown()
}
