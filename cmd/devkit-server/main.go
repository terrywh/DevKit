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
	DefaultConfig.init(fw)
	DefaultOption.init()
	flag.Parse()

	stream.InitTransport(stream.TransportOptions{
		LocalAddress: DefaultOption.addr,
	})
	defer stream.DefaultTransport.Close()

	sc := app.NewServiceController()
	sc.Start(newQuicService())
	sc.WaitForSignal()
	sc.Shutdown()
}
