package k8s

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"git.woa.com/terryhaowu/hybrid-utility/util"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
)


var DefaultTkeClient *tke.Client
var DefaultTkeCredential common.CredentialIface

var oneInitializer *sync.Once = &sync.Once{}
var regexCluster *regexp.Regexp = regexp.MustCompile(`cls-[^\.\s]+`)

type Kubeconfig struct {
	path string // /*/*/{clusterId}.kubeconfig
}

func NewKubeconfig(path string) *Kubeconfig {
	return &Kubeconfig{path}
}

func (c *Kubeconfig) ClusterId() string {
	return regexCluster.FindString(c.path)
}

func (c *Kubeconfig) Exists() bool {
	_, err := os.Stat(c.path)
	return err == nil
}

// Download 下载集群凭证
func (c *Kubeconfig) Download(ctx context.Context) (err error) {
	oneInitializer.Do(func() {
		DefaultTkeCredential = common.NewCredential(util.DefaultConfig.Get("cloud.secret_id", ""), util.DefaultConfig.Get("cloud.secret_key", ""))
		DefaultTkeClient, err = tke.NewClient(DefaultTkeCredential, "ap-guangzhou", profile.NewClientProfile())
		if err != nil {
			DefaultTkeClient = nil
		}
	})
	if DefaultTkeClient != nil { // 直接使用公有云接口
		var rsp *tke.DescribeTKEEdgeExternalKubeconfigResponse
		req := tke.NewDescribeTKEEdgeExternalKubeconfigRequest()
		req.ClusterId = util.StrPtr(c.ClusterId())
		if rsp, err = DefaultTkeClient.DescribeTKEEdgeExternalKubeconfigWithContext(ctx, req); err != nil {
			return
		}
		return os.WriteFile(c.path, []byte(*rsp.Response.Kubeconfig), 0o600)
	} else { // 使用机器人接口
		file := filepath.Base(c.path)
		return util.HttpDownload(ctx, fmt.Sprintf("http://9.221.17.22/robot/cluster/%s", file), c.path)
	}
}
