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
	cert, err := tls.LoadX509KeyPair(DefaultConfig.Get().Client.Certificate,
		DefaultConfig.Get().Client.PrivateKey)
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
		LocalAddress: DefaultConfig.Get().Client.Address,
	})
	defer stream.DefaultTransport.Close()

	sc := app.NewServiceController()
	mgr := stream.NewSessionManager(stream.NewDefaultConnectionProvider(&stream.DialOptions{
		Address:             DefaultConfig.Get().Registry.Address,
		Certificate:         DefaultConfig.Get().Client.Certificate,
		PrivateKey:          DefaultConfig.Get().Client.PrivateKey,
		ApplicationProtocol: "devkit",
	}))
	sc.Start(mgr)
	sc.Start(newServiceHttp(mgr))
	sc.WaitForSignal()
	sc.Close()
}
