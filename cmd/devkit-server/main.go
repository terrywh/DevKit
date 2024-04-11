package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/infra"
	"github.com/terrywh/devkit/stream"
)

func main() {
	fw := infra.NewFileWatcher()
	defer fw.Close()
	DefaultConfig.Init(filepath.Join(app.GetBaseDir(), "etc", "devkit.yaml"))
	fw.Add(DefaultConfig)
	flag.Parse()

	fmt.Println("DeviceID: ", DefaultConfig.Get().DeviceID())

	stream.InitTransport(stream.TransportOptions{
		LocalAddress: DefaultConfig.Get().Server.Address,
	})
	defer stream.DefaultTransport.Close()

	sc := app.NewServiceController()
	sc.Start(fw)
	sc.Start(newQuicService())
	sc.Start(newHttpService())
	if DefaultConfig.Get().Registry.Address != "-" {
		sc.Start(newP2PService())
	}
	sc.WaitForSignal()
	sc.Close()
}
