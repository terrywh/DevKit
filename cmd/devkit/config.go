package main

import (
	"flag"

	"github.com/terrywh/devkit/app"
)

var DefaultConfig *app.Config[ConfigPayload] = &app.Config[ConfigPayload]{}

type ConfigPayloadRegistry struct {
	Address string `yaml:"address"`
}

type ConfigPayloadClient struct {
	Address     string `yaml:"address"`
	Certificate string `yaml:"certificate"`
	PrivateKey  string `yaml:"private_key"`
}

type ConfigPayloadServer struct {
	Address string `yaml:"address"`
}

type ConfigPayload struct {
	Server ConfigPayloadServer `yaml:"server"`
}

func (cp *ConfigPayload) InitFlag(fs *flag.FlagSet) {
	fs.StringVar(&cp.Server.Address, "server.address", "127.0.0.1:18080", "客户及控制服务（本机）")
}
