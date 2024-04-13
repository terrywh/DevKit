package main

import (
	"flag"
	"path/filepath"

	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/infra"
	"github.com/terrywh/devkit/infra/color"
	"github.com/terrywh/devkit/infra/log"
	"github.com/terrywh/devkit/stream"
)

func main() {
	fw := infra.NewFileWatcher()
	defer fw.Close()
	DefaultConfig.Init(filepath.Join(app.GetBaseDir(), "etc", "devkit.yaml"))
	fw.Add(DefaultConfig)
	flag.Parse()

	log.DefaultLogger.SetLevel(log.LevelFromString(DefaultConfig.Get().Log.Level))
	color.Info("DeviceID: ", DefaultConfig.Get().DeviceID(), "\n")

	stream.InitTransport(stream.TransportOptions{
		LocalAddress: DefaultConfig.Get().Server.Address,
	})
	defer stream.DefaultTransport.Close()

	sc := app.NewServiceController()
	defer sc.Close()
	sc.Start(fw)
	sc.Start(newQuicService())
	sc.Start(newHttpService())
	if DefaultConfig.Get().Relay.Address != "-" {
		sc.Start(newP2PService())
	}
	sc.WaitForSignal()
}
