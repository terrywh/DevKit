package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"git.woa.com/terryhaowu/hybrid-utility/k8s"
	"git.woa.com/terryhaowu/hybrid-utility/util"
)

type ClusterServer struct {
	
}

func InitClusterServer(server *http.ServeMux) (svr *ClusterServer) {
	svr = &ClusterServer{}
	server.HandleFunc("/cluster/describe", svr.handleList)
	return
}

func (svr *ClusterServer) handleList(w http.ResponseWriter, r* http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	r.ParseForm()
	req := k8s.Request{
		ClusterID: r.Form.Get("cluster_id"),
		Namespace: r.Form.Get("namespace"),
	}
	if req.ClusterID == "" {
		util.JSONError(w, fmt.Sprint("failed to describe cluster: 'cluster_id' missing"), http.StatusBadRequest)
		return
	}
	desc, err := defaultK8SController.DescribeCluster(ctx, req)
	if err != nil {
		util.JSONError(w, fmt.Sprint("failed to describe cluster: ", desc, err), http.StatusInternalServerError)
		return
	}
	e := json.NewEncoder(w)
	e.Encode(desc)
}
