package entity

type ServerShell struct {
	Server
	ShellId  ShellID  `json:"shell_id,omitempty"`
	ShellCmd []string `json:"shell_cmd"`
	Cols     int      `json:"cols"`
	Rows     int      `json:"rows"`
}
