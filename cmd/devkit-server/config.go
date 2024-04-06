package main

import (
	"log"
	"path/filepath"
	"sync/atomic"

	"github.com/terrywh/devkit/app"
	"github.com/terrywh/devkit/infra"
)

type Config struct {
	path    string
	payload atomic.Pointer[ConfigPayload]
}

func (c *Config) init(fw *infra.FileWatcher) {
	c.path, _ = filepath.Abs(filepath.Join(app.GetBaseDir(), "etc", "devkit.yaml"))
	fw.Add(c)
	c.Reload()
}

func (c *Config) Path() string {
	return c.path
}

func (c *Config) OnChange() {
	c.Reload()
}

func (c *Config) Get() *ConfigPayload {
	return c.payload.Load()
}

func (c *Config) Reload() {
	log.Println("<Config.Reload> ", c.path)
	cp := &ConfigPayload{}
	app.UnmarshalConfig(c.path, cp)
	c.payload.Swap(cp)
}

type ConfigPayload struct {
	Authorize []string `yaml:"authorize"`
}
