package main

import (
	"flag"
	"path/filepath"

	"github.com/terrywh/devkit/app"
)

type Option struct {
	http string
	addr string
	cert string
	pkey string
}

func (o *Option) init() {
	flag.StringVar(&DefaultOption.http, "http", "127.0.0.1:18080", "serve web ui on this address")
	flag.StringVar(&DefaultOption.addr, "bind", "0.0.0.0:18080", "serve QUIC stream on this address")
	flag.StringVar(&DefaultOption.cert, "cert",
		filepath.Join(app.GetBaseDir(), "var/cert/server.crt"), "certificate to use for QUIC (server only)")
	flag.StringVar(&DefaultOption.pkey, "pkey",
		filepath.Join(app.GetBaseDir(), "var/cert/server.key"), "private key to use for QUIC (server only)")
}
