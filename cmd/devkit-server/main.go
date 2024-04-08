package main

import (
	"crypto/tls"
	"flag"
	"log"
	"path/filepath"

	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/infra"
	"github.com/terrywh/devkit/stream"
)

func OutputDeviceID() {
	cert, err := tls.LoadX509KeyPair(DefaultConfig.Get().Server.Certificate,
		DefaultConfig.Get().Server.PrivateKey)
	if err != nil {
		panic("failed to load certificate: " + err.Error())
	}
	log.Println("DeviceID: ", stream.DeviceIDFromCert(cert.Certificate[0]))
}

func main() {
	fw := infra.NewFileWatcher()
	defer fw.Close()
	DefaultConfig.Init(filepath.Join(app.GetBaseDir(), "etc", "devkit.yaml"))
	fw.Add(DefaultConfig)
	flag.Parse()

	OutputDeviceID()

	stream.InitTransport(stream.TransportOptions{
		LocalAddress: DefaultConfig.Get().Server.Address,
	})
	defer stream.DefaultTransport.Close()

	sc := app.NewServiceController()
	sc.Start(fw)
	sc.Start(newQuicService())
	if DefaultConfig.Get().Registry.Address != "-" {
		sc.Start(newP2PService())
	}
	sc.WaitForSignal()
	sc.Close()
}
