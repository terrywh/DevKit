package entity

type RemotePeer struct {
	DeviceID DeviceID `json:"device_id"`
	Address  string   `json:"address"`
}

type RemoteShell struct {
	RemotePeer
	ShellId  ShellID  `json:"shell_id,omitempty"`
	ShellCmd []string `json:"shell_cmd"`
	Cols     int      `json:"cols"`
	Rows     int      `json:"rows"`
}
