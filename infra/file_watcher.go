package infra

import (
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

type FileToWatch interface {
	Path() string
	OnChange()
}

type FileWatcher struct {
	watcher *fsnotify.Watcher

	watch  map[string]FileToWatch
	change map[string]time.Time
}

func (cw *FileWatcher) Add(file FileToWatch) {
	path := file.Path()
	cw.watch[path] = file
	cw.watcher.Add(path)
}

func (cw *FileWatcher) serve() {
	log.Println("<ConfigWatcher.Serve> starting ...")
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
SERVING:
	for {
		select {
		case now := <-ticker.C:
			for name, ct := range cw.change {
				if ct.Before(now) {
					cw.watch[name].OnChange()
					delete(cw.change, name)
				}
			}
		case e := <-cw.watcher.Events:
			log.Println(e.String())
			if e.Op == 0 || e.Name == "" {
				break SERVING
			}
			cw.change[e.Name] = time.Now().Add(2 * time.Second)
		}
	}
	log.Println("<ConfigWatcher.Serve> closed.")
}

func (cw *FileWatcher) Close() (err error) {
	err = cw.watcher.Close()
	return
}

func NewFileWatcher() (fw *FileWatcher) {
	var err error
	fw = &FileWatcher{change: make(map[string]time.Time), watch: make(map[string]FileToWatch)}
	if fw.watcher, err = fsnotify.NewWatcher(); err != nil {
		panic(err)
	}
	go fw.serve()
	return
}
