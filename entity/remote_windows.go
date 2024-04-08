//go:build windows
// +build windows

package entity

func (o *RemoteShell) ApplyDefaults() {
	if len(o.ShellCmd) < 1 {
		o.ShellCmd = []string{"C:\\Windows\\System32\\cmd.exe"}
	}
	if o.Rows < 16 {
		o.Rows = 16
	}
	if o.Cols < 96 {
		o.Cols = 96
	}
}
