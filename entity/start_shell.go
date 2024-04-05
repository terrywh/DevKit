package entity

type StartShell struct {
	DeviceId DeviceID `json:"device_id,omitempty"`
	ShellId  ShellID  `json:"shell_id,omitempty"`
	ShellCmd []string `json:"shell"`
	Cols     int      `json:"cols"`
	Rows     int      `json:"rows"`
}
