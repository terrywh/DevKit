//go:build darwin
// +build darwin

package infra

func (handler *System) Version() (major, minor, build uint32) {
	return 0, 0, 0
}
