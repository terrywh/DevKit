package entity

type DeviceID string

type Server struct {
	DeviceID DeviceID `json:"device_id"`
	Pid      int      `json:"pid"`
	System   string   `json:"system"`
	Arch     string   `json:"arch"`
	Version  string   `json:"version"`
	Address  string   `json:"address"`
}

// Address -> SSH ?
// user:pass@host:port
