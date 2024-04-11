//go:build windows
// +build windows

package main

import "golang.org/x/sys/windows"

func (handler *DeviceHandler) initDeviceVersion() (major, minor, build uint32) {
	return windows.RtlGetNtVersionNumbers()
}
