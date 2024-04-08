package main

import (
	"flag"
	"path/filepath"

	"github.com/terrywh/devkit/app"
)

var DefaultConfig *app.Config[ConfigPayload] = &app.Config[ConfigPayload]{}

type ConfigPayloadClient struct {
	Address string `yaml:"address"`
}

type ConfigPayloadServer struct {
	Address     string `yaml:"address"`
	Certificate string `yaml:"certificate"`
	PrivateKey  string `yaml:"private_key"`
}

type ConfigPayload struct {
	Client ConfigPayloadClient
	Server ConfigPayloadServer
}

func (cp *ConfigPayload) InitFlag() {
	flag.StringVar(&cp.Client.Address, "client.address", "127.0.0.1:18080", "serve web ui on this address")
	flag.StringVar(&cp.Server.Address, "server.address", "0.0.0.0:18080", "serve QUIC stream on this address")
	flag.StringVar(&cp.Server.Certificate, "server.certificate",
		filepath.Join(app.GetBaseDir(), "var/cert/server.crt"), "certificate to use for QUIC (server only)")
	flag.StringVar(&cp.Server.PrivateKey, "server.private_key",
		filepath.Join(app.GetBaseDir(), "var/cert/server.key"), "private key to use for QUIC (server only)")

}
