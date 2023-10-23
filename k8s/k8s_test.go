package k8s

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
)

func TestInstall(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	bin := filepath.Join(filepath.Dir(filepath.Dir(filename)), "bin")
	path := filepath.Join(bin, "kubectl_darwin_arm64")

	
	cmd := NewKubectl(path)
	
	t.Log(cmd.Version(context.Background()))
	t.Log(cmd.Upgrade(context.Background()))
	t.Log(cmd.Version(context.Background()))
}