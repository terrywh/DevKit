package entity

type StartShell struct {
	DeviceId DeviceID `json:"device_id,omitempty"`
	ShellId  ShellID  `json:"shell_id,omitempty"`
	ShellCmd []string `json:"shell"`
	Cols     int      `json:"cols"`
	Rows     int      `json:"rows"`
}

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
