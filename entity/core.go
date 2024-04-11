package entity

type DeviceID string
type ShellID string

type Server struct {
	DeviceID DeviceID `json:"device_id"`
	Address  string   `json:"address"`
	System   string   `json:"system"`
	Arch     string   `json:"arch"`
	Version  string   `json:"version"`
}

type ServerShell struct {
	Server
	ShellId  ShellID  `json:"shell_id,omitempty"`
	ShellCmd []string `json:"shell_cmd"`
	Cols     int      `json:"cols"`
	Rows     int      `json:"rows"`
}
