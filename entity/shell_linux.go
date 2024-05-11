//go:build linux
// +build linux

package entity

func (o *Shell) ApplyDefaults() {
	if len(o.Cmd) < 1 {
		o.Cmd = []string{"bash"}
	}
	if o.Row < 16 {
		o.Row = 16
	}
	if o.Col < 96 {
		o.Col = 96
	}
}
