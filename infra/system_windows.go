//go:build windows
// +build windows

package infra

import "golang.org/x/sys/windows"

func (system *System) Version() (major, minor, build uint32) {
	return windows.RtlGetNtVersionNumbers()
}
