//go:build !windows
// +build !windows

package main

func (handler *DeviceHandler) initDeviceVersion() (major, minor, build uint32) {
	return 0, 0, 0
}
