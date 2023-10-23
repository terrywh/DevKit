package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"path/filepath"
	"runtime"
	"time"

	"github.com/getlantern/systray"
	"github.com/terrywh/devkit/k8s"
	"github.com/terrywh/devkit/ssh"
	"github.com/terrywh/devkit/util"
)

var defaultSSHController * ssh.Controller
var defaultK8SController * k8s.Controller

func devkitAppServer() {
	util.OnInit(defaultConfig)
	defaultSSHController = ssh.NewController() 
	defaultK8SController = k8s.NewController()

	server := http.NewServeMux()
	
	apiServer := InitAppServer(defaultConfig.Local.Root, server)
	InitBashServer(server)
	InitClusterServer(server)

	systray.Run(apiServer.onReady, apiServer.onExit)
}

func upgradeKubectl(ctx context.Context) {
	cmd := k8s.NewKubectl(filepath.Join(defaultConfig.Local.Root,
		fmt.Sprintf("bin/kubectl_%s_%s", runtime.GOOS, runtime.GOARCH)))
	
	ctx, cancel := context.WithTimeout(ctx, 30 * time.Second)
	defer cancel()
	if v, err := cmd.Upgrade(ctx); err == nil {
		log.Println("done", v)
	} else {
		log.Println("fail", err)
	}
}

func main() {
	ctx := context.Background()
	defaultConfig.Init()

	flagUpgrade := flag.String("upgrade", "", "升级安装 kubectl 依赖组件")
	flagHelp := flag.Bool("help", false, "查看帮助信息")
	flag.Parse()
	if *flagHelp {
		flag.PrintDefaults()
	} else if *flagUpgrade == "kubectl" {
		upgradeKubectl(ctx)
	} else if *flagUpgrade == "trzsz" {
		log.Println("暂不支持")
	} else {
		devkitAppServer()
	}
}
