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
		LocalAddress: DefaultConfig.Get().Client.Address,
	})
	defer stream.DefaultTransport.Close()

	sc := app.NewServiceController()
	opts := &stream.DialOptions{
		Address:             DefaultConfig.Get().Registry.Address,
		Certificate:         DefaultConfig.Get().Client.Certificate,
		PrivateKey:          DefaultConfig.Get().Client.PrivateKey,
		ApplicationProtocol: "devkit",
	}
	mux := stream.NewServeMux()
	mgr := stream.NewSessionManager(&stream.SessionManagerOptions{
		DialOptions: *opts,
		Resolver:    newResolver(opts),
		Handler: &stream.DefaultConnectionHandler{
			Tracker: stream.NewDefaultConnectionTracker(),
			Handler: mux,
		},
	})
	initFileHandler(mgr, mux)
	sc.Start(mgr)
	sc.Start(newServiceHttp(mgr, mux))
	sc.WaitForSignal()
	sc.Close()
}
