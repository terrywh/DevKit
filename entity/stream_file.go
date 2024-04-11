package entity

type File struct {
	Path string `json:"path"`
	Size int64  `json:"size"`
	Perm uint32 `json:"perm"`
}

type FilePull struct {
	File
	DeviceID DeviceID `json:"device_id"`
}

type FilePush struct {
	File
	Override bool `json:"override"`
}

type ServerStreamFilePull struct {
	FilePull
	Pid int `json:"pid"`
}
