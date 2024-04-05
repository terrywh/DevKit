package main

import (
	"flag"
	"net/http"
	_ "net/http/pprof"

	"github.com/getlantern/systray"
	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/k8s"
	"github.com/terrywh/devkit/ssh"
)

var defaultSSHController *ssh.Controller
var defaultK8SController *k8s.Controller

func devkitAppServer() {
	app.InitConfigWatcher("devkit", defaultConfig)
	defaultSSHController = ssh.NewController()
	defaultK8SController = k8s.NewController()

	server := http.NewServeMux()

	apiServer := InitAppServer(defaultConfig.Local.Root, server)
	InitBashServer(server)
	InitClusterServer(server)

	systray.Run(apiServer.onReady, apiServer.onExit)
}

func main() {
	// ctx := context.Background()
	defaultConfig.Init()

	flagHelp := flag.Bool("help", false, "查看帮助信息")
	flag.Parse()
	if *flagHelp {
		flag.PrintDefaults()
	} else {
		devkitAppServer()
	}
}
