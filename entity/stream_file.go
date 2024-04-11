package entity

type StreamFile struct {
	Path string `json:"path"`
	Size int64  `json:"size"`
	Perm uint32 `json:"perm"`
}

type StreamFilePull struct {
	StreamFile
	DeviceID DeviceID `json:"device_id"`
}

type StreamFilePush struct {
	StreamFile
	Override bool `json:"override"`
}

type ServerStreamFilePull struct {
	StreamFilePull
	Pid int `json:"pid"`
}
