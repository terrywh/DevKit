package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/getlantern/systray"
)

type AppServer struct {
	server *http.ServeMux
	root string
}


func InitAppServer(server *http.ServeMux) (api *AppServer) {
	api = &AppServer{ server, "./" }
	server.HandleFunc("/app/login", api.handleLogin)
	server.HandleFunc("/app/quit", api.handleQuit)
	return
}

func (api *AppServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (api *AppServer) handleQuit(w http.ResponseWriter, r *http.Request) {
	systray.Quit()
	w.WriteHeader(http.StatusOK)
}

func (api *AppServer) detectRootDir() {
	detect := func (path string) bool {
		if _, err := os.Stat(filepath.Join(path, "public")); err != nil {
			return false
		}
		return true
	}
	bin, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	if detect(bin) {
		api.root = bin
		return
	}
	if detect(filepath.Dir(bin)) {
		api.root = filepath.Dir(bin)
		return
	}
	_, filename, _, _ := runtime.Caller(0)
	api.root = filepath.Dir(filepath.Dir(filepath.Dir(filename)))
}

func (api *AppServer) onReady() {
	api.detectRootDir()
	log.Println("onReady, start server (", api.root, ") :8080 ...")
	go http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		staticMOD := http.FileServer(http.Dir(filepath.Join(api.root, "node_modules")))
		staticWWW := http.FileServer(http.Dir(filepath.Join(api.root, "public")))
	
		var found bool
		if handler, pattern := api.server.Handler(r); pattern != "" {
			handler.ServeHTTP(w, r)
		} else if r.URL.Path, found = strings.CutPrefix(r.URL.Path, "/node_modules"); found && r.URL.Path[0] == '/' {
			staticMOD.ServeHTTP(w, r)
		} else if strings.HasPrefix(r.URL.Path, "/debug/pprof") {
			http.DefaultServeMux.ServeHTTP(w, r)
		} else {
			staticWWW.ServeHTTP(w, r)
		}
	}))
	
	// systray.SetIcon()
	systray.SetTitle("devkit")
	menuItem := systray.AddMenuItem("Quit", "Quit Wemeet-Hybrid")
	go func() {
		<- menuItem.ClickedCh
		systray.Quit()
	} () 
}

func (api *AppServer)  onExit() {
	log.Println("onExit.")
}


