package entity

type ShellID string

type Shell struct {
	ID  ShellID  `json:"id,omitempty"`
	Cmd []string `json:"cmd"`
	Col int      `json:"col"`
	Row int      `json:"row"`
}

type ServerShell struct {
	Server  Server  `json:"server,omitempty"`
	Cluster Cluster `json:"cluster,omitempty"`
	Shell   Shell   `json:"shell"`
}
