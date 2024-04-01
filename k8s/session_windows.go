//go:build windows
// +build windows

package k8s

import (
	"context"
	"fmt"

	"github.com/UserExistsError/conpty"
)

func StartSession(ctx context.Context, s *Session) (pty Pseudo, err error) {
	return conpty.Start(
		fmt.Sprintf("%s --kubeconfig %s exec -n %s -it %s -- bash", s.path, s.conf, s.Req.Namespace, s.Req.Pod),
		conpty.ConPtyDimensions(s.Req.Cols, s.Req.Rows))
}
