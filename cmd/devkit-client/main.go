package main

import (
	"flag"
	"fmt"
	"net"
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

	_, port, _ := net.SplitHostPort(DefaultConfig.Get().Client.Address)
	stream.InitTransport(stream.TransportOptions{
		LocalAddress: fmt.Sprintf("0.0.0.0:%s", port),
	})
	defer stream.DefaultTransport.Close()

	sc := app.NewServiceController()
	defer sc.Close()
	// go func() {
	opts := &stream.DialOptions{
		Address:             DefaultConfig.Get().Relay.Address,
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
	sc.Start(newServiceHttp(mgr))
	sc.WaitForSignal()
	// }()
	// 必须在主线程运行
	// wv := newServiceWebview()
	// defer wv.Close()
	// wv.Serve(context.Background())
}
