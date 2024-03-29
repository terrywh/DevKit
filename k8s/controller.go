package k8s

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/terrywh/devkit/util"
)

type Controller struct {
	Executable string
	ConfigDir  string
}

func NewController() *Controller {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".kube") // 配置文件存储路径(与 kubectl 保持一致)
	os.Mkdir(path, 0o755)
	var ext string
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}
	return &Controller{
		Executable: fmt.Sprintf("./bin/kubectl_%s_%s%s",
			runtime.GOOS, runtime.GOARCH, ext),
		ConfigDir: path,
	}
}

func (c *Controller) config(ctx context.Context, clusterID string) (path string, err error) {
	path = filepath.Join(c.ConfigDir, fmt.Sprintf("%s.kubeconfig", clusterID))
	conf := NewKubeconfig(path)
	if !conf.Exists() {
		if err = conf.Download(ctx); err != nil {
			return
		}
	}
	return
}

func (c *Controller) CreateShell(ctx context.Context, req Request) (session *Session, err error) {
	var conf string
	if conf, err = c.config(ctx, req.ClusterID); err != nil {
		return
	}
	session = &Session{
		Req:  req,
		path: c.Executable,
		conf: conf,
	}
	return session, nil
}

type ListItem struct {
	Items []util.JSONObject `json:"items"`
}

func (c *Controller) buildPort(ports []interface{}) (r string) {
	for _, port := range ports {
		obj := util.JSONObject(port.(map[string]interface{}))
		from := obj.GetString("nodePort")
		if from == "" {
			from = obj.GetString("port")
		}
		r += fmt.Sprintf("%s => %s ", from, obj.GetString("targetPort"))
	}
	return
}

func (c *Controller) buildStat(conditions []interface{}) string {
	for _, stat := range conditions {
		obj := util.JSONObject(stat.(map[string]interface{}))
		if obj.GetString("type") == "Ready" {
			if status := obj.GetString("status"); status == "True" {
				return "Ready"
			} else if status == "False" {
				return "Unhealthy"
			} else {
				return "Unknown"
			}
		}
	}
	return "Unknown"
}

func (c *Controller) buildAddr(addrs []interface{}) string {
	for _, addr := range addrs {
		obj := util.JSONObject(addr.(map[string]interface{}))
		if obj.GetString("type") == "InternalIP" {
			return obj.GetString("address")
		}
	}
	return ""
}

func (c *Controller) DescribeCluster(ctx context.Context, req Request) (desc *Description, err error) {
	var conf string
	if conf, err = c.config(ctx, req.ClusterID); err != nil {
		return
	}
	log.Println(c.Executable, "--kubeconfig", conf, "-n", req.Namespace, "-o", "json", "get", "pod,svc,node")
	proc := exec.CommandContext(ctx, c.Executable, "--kubeconfig", conf, "-n", req.Namespace, "-o", "json", "get", "node,pod,svc")
	var data []byte
	if data, err = proc.Output(); err != nil {
		return
	}
	desc = &Description{
		ClusterID: req.ClusterID,
		Namespace: req.Namespace,
	}
	var store ListItem
	d := json.NewDecoder(bytes.NewBuffer(data))
	d.Decode(&store)
	for _, object := range store.Items {
		if kind := object.Get("kind"); kind == "Service" {

			if _, ok := object.Get("spec.ports").([]interface{}); !ok {
				log.Println(object.Get("spec.ports"))
			}

			svc := Svc{
				Name:      object.GetString("metadata.name"),
				Type:      object.GetString("spec.type"),
				ClusterIP: object.GetString("spec.clusterIP"),
				// Port:      c.buildPort(),
			}
			if ports, ok := object.Get("spec.ports").([]interface{}); ok {
				svc.Port = c.buildPort(ports)
			}
			desc.Svc = append(desc.Svc, svc)
		} else if kind == "Pod" {
			desc.Pod = append(desc.Pod, Pod{
				Name:   object.GetString("metadata.name"),
				Node:   object.GetString("status.hostIP"),
				IP:     object.GetString("status.podIP"),
				Status: object.GetString("status.phase"),
			})
		} else if kind == "Node" {
			desc.Node = append(desc.Node, Node{
				IP:     c.buildAddr(object.Get("status.addresses").([]interface{})),
				Status: c.buildStat(object.Get("status.conditions").([]interface{})),
			})
		}
	}
	return
}

func (c *Controller) Cleanup(ctx context.Context) error {
	return nil
}
