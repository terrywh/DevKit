package main

import (
	"flag"
	"path/filepath"

	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/util"
)

var DefaultConfig *app.Config[ConfigPayload] = &app.Config[ConfigPayload]{}

type ConfigPayloadRegistry struct {
	Address string `yaml:"address"`
}

type ConfigPayloadClient struct {
	Address string `yaml:"address"`
}

type ConfigPayloadServer struct {
	Address     string          `yaml:"address"`
	Certificate string          `yaml:"certificate"`
	PrivateKey  string          `yaml:"private_key"`
	Authorized  util.StringList `yaml:"authorized"`
}

type ConfigPayload struct {
	Registry ConfigPayloadRegistry `yaml:"registry"`
	Client   ConfigPayloadClient   `yaml:"client"`
	Server   ConfigPayloadServer   `yaml:"server"`
}

func (cp *ConfigPayload) InitFlag() {
	flag.StringVar(&cp.Registry.Address, "registry.address", "42.193.117.122:18080", "注册呼叫服务")
	flag.StringVar(&cp.Client.Address, "client.address", "127.0.0.1:18080", "客户及控制服务")
	flag.StringVar(&cp.Server.Address, "server.address", "0.0.0.0:18080", "服务监听")
	flag.StringVar(&cp.Server.Certificate, "server.certificate",
		filepath.Join(app.GetBaseDir(), "var/cert/server.crt"), "服务证书公钥")
	flag.StringVar(&cp.Server.PrivateKey, "server.private_key",
		filepath.Join(app.GetBaseDir(), "var/cert/server.key"), "服务证书私钥")
	flag.Var(&cp.Server.Authorized, "server.authorized", "认证客户端")

}
