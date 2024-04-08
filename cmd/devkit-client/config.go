package main

import (
	"flag"
	"path/filepath"

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
}

type ConfigPayload struct {
	Registry ConfigPayloadRegistry `yaml:"registry"`
	Client   ConfigPayloadClient   `yaml:"client"`
	Server   ConfigPayloadServer   `yaml:"server"`
}

func (cp *ConfigPayload) InitFlag() {
	flag.StringVar(&cp.Registry.Address, "registry.address", "42.193.117.122:18080", "注册呼叫服务")
	flag.StringVar(&cp.Client.Address, "client.address", "127.0.0.1:18080", "客户及控制服务")
	flag.StringVar(&cp.Client.Certificate, "client.certificate",
		filepath.Join(app.GetBaseDir(), "var/cert/client.crt"), "连接认证证书公钥")
	flag.StringVar(&cp.Client.PrivateKey, "client.private_key",
		filepath.Join(app.GetBaseDir(), "var/cert/client.key"), "连接认证证书私钥")
}
