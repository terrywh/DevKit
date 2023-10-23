package main

import (
	"os"
	"path/filepath"
	"runtime"
)

// DefaultConfig ...
type DefaultConfig struct {
	Local struct {
		Root string
	}
	Cloud struct {
		SecretID  string `yaml:"secret_id"`
		SecretKey string `yaml:"secret_key"`
	} `yaml:"cloud"`
}

// Get ...
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

// 
func (c *DefaultConfig) Init() {
	c.initLocalRoot()
}

func (c *DefaultConfig) initLocalRoot() {
	detect := func (path string) bool {
		if _, err := os.Stat(filepath.Join(path, "public")); err != nil {
			return false
		}
		return true
	}
	bin, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	if detect(bin) {
		c.Local.Root = bin
		return
	}
	if detect(filepath.Dir(bin)) {
		c.Local.Root = filepath.Dir(bin)
		return
	}
	_, filename, _, _ := runtime.Caller(0)
	c.Local.Root = filepath.Dir(filepath.Dir(filepath.Dir(filename)))
}
