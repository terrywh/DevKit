package k8s

import (
	"context"
	"testing"
)

func TestDescribe(t *testing.T) {
	c := NewController()
	if desc, err := c.DescribeCluster(context.Background(), Request{
		ClusterID: "cls-s0d109ge",
		Namespace: "wemeet",
	}); err != nil {
		t.Error("failed to describe cluster: ", desc, err)
	}
}