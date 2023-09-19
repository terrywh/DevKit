package k8s

import (
	"context"
	"fmt"
	"runtime"

	"github.com/terrywh/dev-kit/util"
)

// Request ...
type Request struct {
	ClusterID string `json:"cluster_id"`
	Namespace string `json:"namespace"`
	Pod string `json:"pod"`
	Command string `json:"command"`

	Rows int `json:"rows"`
	Cols int `json:"cols"`
}
type Description struct {
	ClusterID string `json:"cluster_id"`
	Namespace string `json:"namespace"`
	Svc []Svc `json:"svc"`
	Pod []Pod `json:"pod"`
	Node []Node `json:"node"`
}
type Svc struct {
	Name string `json:"name"`
	Type string `json:"type"`
	ClusterIP string `json:"cluster_ip"`
	Port string `json:"port"`
}
type Pod struct {
	Name string `json:"name"`
	Node string `json:"node"`
	IP string `json:"ip"`
	Status string `json:"status"`
}
type Node struct {
	IP string `json:"ip"`
	Status string `json:"status"`
}

// KubectlClientVersion ...
type KubectlClientVersion struct {
	ClientVersion struct {
		GitVersion string `json:"gitVersion"`
	} `json:"clientVersion"`
}

// GetKubectlStableVersion 获取对应此命令行工具的稳定版本
func GetKubectlStableVersion(ctx context.Context) (version string) {
	version, _ = util.HttpGet(ctx, "https://dl.k8s.io/release/stable.txt")
	return version
}

// Kubectl 命令行工具的简单封装
type Kubectl struct {
	Path string
}

// NewKubectl ...
func NewKubectl(path string) (cmd *Kubectl) {
	return &Kubectl{ path }
}

// Version ...
func (cmd *Kubectl) Version(ctx context.Context) string {
	var version KubectlClientVersion
	util.FileExecuteJSON(ctx, &version, cmd.Path, "version", "--client=true", "--output=json")
	return version.ClientVersion.GitVersion
}
// Upgrade ...
func (cmd *Kubectl) Upgrade(ctx context.Context) (string, error) {
	stablev := GetKubectlStableVersion(ctx)
	version := cmd.Version(ctx)
	if version != stablev {
		err := cmd.Install(ctx, stablev)
		return stablev, err
	}
	return version, nil
}

func (cmd *Kubectl) Install(ctx context.Context, version string) (err error) {
	// https://dl.k8s.io/release/v1.27.3/bin/darwin/arm64/kubectl
	url := fmt.Sprintf("https://dl.k8s.io/release/%s/bin/%s/%s/kubectl", version, runtime.GOOS, runtime.GOARCH)
	err = util.HttpDownload(ctx, url, cmd.Path)
	util.FileAddPerm(cmd.Path, 0o111)
	return
}

