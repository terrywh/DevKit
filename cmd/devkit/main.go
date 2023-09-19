package main

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/getlantern/systray"
	"github.com/terrywh/dev-kit/k8s"
	"github.com/terrywh/dev-kit/ssh"
	"github.com/terrywh/dev-kit/util"
)

type DefaultConfig struct {
	Cloud struct {
		SecretID  string `yaml:"secret_id"`
		SecretKey string `yaml:"secret_key"`
	} `yaml:"cloud"`
}

func (c *DefaultConfig) Get(field string, defval string) (value string) {
	switch field {
	case "cloud.secret_id":
		return c.Cloud.SecretID
	case "cloud.secret_key":
		return c.Cloud.SecretKey
	default:
		return defval
	}
}

var defaultConfig *DefaultConfig = &DefaultConfig{}
var defaultSSHController * ssh.Controller
var defaultK8SController * k8s.Controller


func main() {
	util.OnInit(defaultConfig)
	defaultSSHController = ssh.NewController() 
	defaultK8SController = k8s.NewController()

	server := http.NewServeMux()
	
	apiServer := InitAppServer(server)
	InitBashServer(server)
	InitClusterServer(server)

	systray.Run(apiServer.onReady, apiServer.onExit)
}
