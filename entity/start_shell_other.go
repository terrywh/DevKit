//go:build !windows
// +build !windows

package entity

func (o *StartShell) ApplyDefaults() {
	if len(o.ShellCmd) < 1 {
		o.ShellCmd = []string{"bash"}
	}
	if o.Rows < 16 {
		o.Rows = 16
	}
	if o.Cols < 96 {
		o.Cols = 96
	}
}
