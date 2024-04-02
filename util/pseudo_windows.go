//go:build windows
// +build windows

package util

import (
	"context"
	"strings"

	"github.com/UserExistsError/conpty"
)

func StartPty(ctx context.Context, rows, cols int, cmd string, args ...string) (pty Pseudo, err error) {
	return conpty.Start(strings.Join([]string{cmd, strings.Join(args, " ")}, " "),
		conpty.ConPtyDimensions(cols, rows))
}
