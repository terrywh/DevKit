package entity

type File struct {
	Path string `json:"path"`
	Size int64  `json:"size,omitempty"`
	Perm uint32 `json:"perm,omitempty"`
}

type StreamFile struct {
	Source  File              `json:"source,omitempty"`
	Target  File              `json:"target,omitempty"`
	Options StreamFileOptions `json:"options"`
}

type StreamFileOptions struct {
	Override bool `json:"override,omitempty"`
}
