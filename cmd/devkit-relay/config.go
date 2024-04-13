package main

import (
	"flag"
	"path/filepath"

	"github.com/terrywh/devkit/app"
)

var DefaultConfig *app.Config[ConfigPayload] = &app.Config[ConfigPayload]{}

type ConfigPayloadServer struct {
	Address     string `yaml:"address"`
	Certificate string `yaml:"certificate"`
	PrivateKey  string `yaml:"private_key"`
}

type ConfigPayloadLog struct {
	Level string `yaml:"level"`
}
type ConfigPayload struct {
	Server ConfigPayloadServer `yaml:"server"`
	Log    ConfigPayloadLog    `yaml:"log"`
}

func (cp *ConfigPayload) InitFlag() {
	flag.StringVar(&cp.Log.Level, "log.level", "warn", "日志级别")
	flag.StringVar(&cp.Server.Address, "server.address", "0.0.0.0:18080", "注册呼叫服务")
	flag.StringVar(&cp.Server.Certificate, "server.certificate",
		filepath.Join(app.GetBaseDir(), "var/cert/server.crt"), "服务证书公钥")
	flag.StringVar(&cp.Server.PrivateKey, "server.private_key",
		filepath.Join(app.GetBaseDir(), "var/cert/server.key"), "服务证书私钥")

}
