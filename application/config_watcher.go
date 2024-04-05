package application

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/terrywh/devkit/util"
	"gopkg.in/yaml.v3"
)

type ConfigWatcher struct {
	name    string
	base    string
	path    string
	dest    interface{}
	watcher *fsnotify.Watcher
	cancel  context.CancelFunc
}

func (cw *ConfigWatcher) initBaseDir() {
	bin, _ := os.Executable()
	if cw.name == "" {
		cw.name = filepath.Base(DefaultConfigWatcher.name)
		cw.name, _ = strings.CutSuffix(cw.name, filepath.Ext(cw.name))
	}
	DefaultConfigWatcher.base, _ = filepath.Abs(bin)
	for i := 0; i < 3; i++ {
		DefaultConfigWatcher.base = filepath.Dir(DefaultConfigWatcher.base)
		if _, err := os.Stat(filepath.Join(DefaultConfigWatcher.base, "bin")); os.IsNotExist(err) {
			continue
		}
		break
	}
}

func (cw *ConfigWatcher) initConfig() {
	cw.path = filepath.Join(cw.base, "etc", fmt.Sprintf("%s.yaml", cw.name))
	cw.watcher.Add(cw.path)
	ParseConfig(cw.path, cw.dest)
}

func (cw *ConfigWatcher) BaseDir() string {
	return cw.base
}

func (cw *ConfigWatcher) Serve(ctx context.Context) {
	forParseConfig := util.NewThrottle(3 * time.Second)
SERVING:
	for {
		select {
		case <-ctx.Done():
			break SERVING
		case e := <-cw.watcher.Events:
			switch e.Op {
			case fsnotify.Write:
				if e.Name == cw.path {
					forParseConfig.Do(func() {
						ParseConfig(cw.path, cw.dest)
					})
				}
			}
		}
	}
}

func (cw *ConfigWatcher) Close() (err error) {
	err = cw.watcher.Close()
	cw.cancel()
	return
}

var DefaultConfigWatcher *ConfigWatcher

func InitConfigWatcher(name string, conf interface{}) (err error) {
	DefaultConfigWatcher = &ConfigWatcher{name: name, dest: conf}
	DefaultConfigWatcher.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	go DefaultConfigWatcher.Serve(ctx)
	DefaultConfigWatcher.cancel = cancel

	DefaultConfigWatcher.initBaseDir()
	DefaultConfigWatcher.initConfig()

	return
}

func ParseConfig(path string, v interface{}) (err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	err = yaml.NewDecoder(file).Decode(v)
	log.Println("<ParseConfig> load config on: ", path)
	return
}
