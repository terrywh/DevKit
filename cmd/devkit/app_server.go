package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/getlantern/systray"
)

type AppServer struct {
	server *http.ServeMux
	root   string
}

// InitAppServer ...
func InitAppServer(root string, server *http.ServeMux) (api *AppServer) {
	api = &AppServer{server, root}
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

func (api *AppServer) onReady() {
	log.Println("onReady, start server (", api.root, ") :8080 ...")
	go http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		staticMOD := http.FileServer(http.Dir(filepath.Join(api.root, "node_modules")))
		staticWWW := FileServer{Root: filepath.Join(api.root, "public"), Mime: map[string]string{".svelte": "text/javascript"}}

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
	icon, _ := os.ReadFile("var/icon.ico")
	systray.SetIcon(icon)
	systray.SetTooltip("devkit")
	// systray.SetTitle("devkit")
	menuItem := systray.AddMenuItem("Quit", "Quit devkit")
	go func() {
		<-menuItem.ClickedCh
		systray.Quit()
	}()
}

func (api *AppServer) onExit() {
	log.Println("onExit.")
}
