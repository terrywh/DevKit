package main

import (
	"flag"
	"path/filepath"

	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/infra"
	"github.com/terrywh/devkit/stream"
)

func main() {
	DefaultConfig.Init(filepath.Join(app.GetBaseDir(), "etc", "devkit.yaml"))

	flag.Parse()

	stream.InitTransport(stream.TransportOptions{
		LocalAddress: DefaultConfig.Get().Server.Address,
	})
	defer stream.DefaultTransport.Close()

	sc := app.NewServiceController()
	defer sc.Close()

	fw := infra.NewFileWatcher()
	fw.Add(DefaultConfig)
	sc.Start(fw)
	sc.Start(newQuicService())

	sc.WaitForSignal()
}
