package main

import (
	"crypto/tls"
	"flag"
	"path/filepath"

	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/stream"
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
	deviceID entity.DeviceID
	cert     tls.Certificate
}

func (cp *ConfigPayload) InitFlag() {
	flag.StringVar(&cp.Registry.Address, "registry.address", "42.193.117.122:18080", "注册呼叫服务")
	flag.StringVar(&cp.Client.Address, "client.address", "0.0.0.0:18080", "客户及控制服务")
	flag.StringVar(&cp.Client.Certificate, "client.certificate",
		filepath.Join(app.GetBaseDir(), "var/cert/client.crt"), "服务证书公钥")
	flag.StringVar(&cp.Client.PrivateKey, "client.private_key",
		filepath.Join(app.GetBaseDir(), "var/cert/client.key"), "服务证书私钥")
}

func (cp *ConfigPayload) Certificate() tls.Certificate {
	if len(cp.cert.Certificate) > 0 {
		return cp.cert
	}
	var err error
	cp.cert, err = tls.LoadX509KeyPair(cp.Client.Certificate, cp.Client.PrivateKey)
	if err != nil {
		panic("failed to load certificate: " + err.Error())
	}
	return cp.cert
}

func (cp *ConfigPayload) DeviceID() entity.DeviceID {
	if len(cp.deviceID) > 0 {
		return cp.deviceID
	}
	cp.deviceID = stream.DeviceIDFromCert(cp.Certificate().Certificate[0])
	return cp.deviceID
}
