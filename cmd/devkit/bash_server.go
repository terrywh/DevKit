package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/terrywh/devkit/util"
	"golang.org/x/net/websocket"
)

type BashServer struct {
	ws websocket.Handler
	c  *BashController
}

func InitBashServer(server *http.ServeMux) (api *BashServer) {
	api = &BashServer{}
	api.ws = websocket.Handler(api.handleStreamSocket)
	api.c = NewBashController()

	server.HandleFunc("/bash/create", api.handleCreate)
	server.HandleFunc("/bash/stream", api.handleStream)
	server.HandleFunc("/bash/resize", api.handleResize)
	server.HandleFunc("/bash/config", api.handleConfig)
	return
}

func (svr *BashServer) handleCreate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if _, err := svr.c.FetchShell(ctx, r); err != nil {
		util.JSONError(w, fmt.Sprint("failed to create shell: ", err), http.StatusInternalServerError)
		return
	} else {
		fmt.Fprintf(w, `{"key":"%s"}`, r.URL.Query().Get("key"))
	}
}

func (svr *BashServer) handleStreamSocket(conn *websocket.Conn) {
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	shell, err := svr.c.FetchShell(ctx, conn.Request())
	// bash, err := defaultSshController.CreateShell(context.Background(), ssh.Request{
	// 	Host: "csig.mnet2.com",
	// 	Port: 36000,
	// 	User: "terryhaowu",
	// 	Command: "tkex-login -cls cls-iduhpo4b -n ns-prj62vsz-1512194-test -p wemeet-hybrid-proxy-server-ci-38-env-59180028-0 -c hybrid-proxy-server -b /bin/bash",
	// })
	// bash, err := defaultK8sController.CreateShell(context.Background(), k8s.Request{
	// 	ClusterId: "cls-s0d109ge",
	// 	Namespace: "wemeet",
	// 	Pod: "wemeet-hybrid-proxy-agent-0",
	// })
	if err != nil {
		log.Println("failed to fetch ssh shell: ", shell, err)
		return
	}
	defer shell.Close()

	log.Println("<bash/stream> start ...")
	bash := &BashShell {source: BashSocket{conn}, server: shell}
	bash.Serve(context.Background())
	log.Println("<bash/stream> close.")
}

func (svr *BashServer) handleStream(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !r.Form.Has("key") {
		util.JSONError(w, "failed to create shell stream: missing key", http.StatusBadRequest)
		return
	}
	svr.ws.ServeHTTP(w, r)
}

func (svr *BashServer) handleResize(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r.ParseForm()
	if shell, err := svr.c.FetchShell(ctx, r); err != nil {
		util.JSONError(w, fmt.Sprint("failed to fetch shell: ", err), http.StatusPreconditionFailed)
		return
	} else {
		decoder := json.NewDecoder(r.Body)
		var req ResizeRequestImpl
		decoder.Decode(&req)
		shell.Resize(req.Rows(), req.Cols())
		fmt.Fprintf(w, `{"r":%d}`, time.Now().Unix())
	}
}

func (svr *BashServer) handleConfig(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r.ParseForm()
	shell, err := svr.c.FetchShell(ctx, r)
	if err != nil {
		util.JSONError(w, fmt.Sprint("failed to fetch shell: ", err), http.StatusPreconditionFailed)
		return
	}
	
	setup := &BashSetup{ shell }
	setup.Serve(context.Background())
	fmt.Fprintf(w, `{"r":%d}`, time.Now().Unix())
}
