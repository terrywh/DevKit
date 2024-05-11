//go:build !windows
// +build !windows

package infra

import (
	"os/exec"
	"regexp"

	"github.com/terrywh/devkit/util"
)

var extractVersionPattern = regexp.MustCompile(`(\d+)\.(\d+)(\.(\d+))?`)

func (handler *System) Version() (major, minor, build uint32) {
	cmd := exec.Command("uname", "-r")
	if out, err := cmd.Output(); err != nil {
		return
	} else {
		tmp := extractVersionPattern.FindStringSubmatch(string(out))
		major = util.ToInteger[uint32](tmp[1])
		minor = util.ToInteger[uint32](tmp[2])
		build = util.ToInteger[uint32](tmp[4])
		return
	}
}
