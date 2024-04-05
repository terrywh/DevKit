package main

import (
	"flag"

	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/infra"
	"github.com/terrywh/devkit/stream"
)

var DefaultOption Option
var DefaultConfig Config

func main() {
	fw := infra.NewFileWatcher()
	defer fw.Close()
	DefaultOption.init()
	DefaultConfig.init(fw)
	flag.Parse()

	stream.InitTransport(stream.TransportOptions{
		LocalAddress: DefaultOption.addr,
	})
	defer stream.DefaultTransport.Close()

	stream.InitSessionManager(stream.NewDirectProvider(stream.DialOptions{
		Certificate:         DefaultOption.cert,
		PrivateKey:          DefaultOption.pkey,
		ApplicationProtocol: "devkit",
	}))
	defer stream.DefaultSessionManager.Close()

	sc := app.NewServiceController()
	sc.Start(newServiceHttp())
	sc.WaitForSignal()
	sc.Shutdown()
}
