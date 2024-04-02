package stream

import (
	"encoding/json"
	"io"
)

type Request interface {
	ApplyDefaults()
	WriteTo(w io.Writer) error
}

type StreamFileReq struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Perm uint32 `json:"perm"`
}

func (req *StreamFileReq) ApplyDefaults() {}

func (req *StreamFileReq) WriteTo(w io.Writer) (err error) {
	io.WriteString(w, "StreamFile:")
	e := json.NewEncoder(w)
	err = e.Encode(req)
	return
}

type OpenShellReq struct {
	Shell string `json:"shell"`
	Cols  uint16 `json:"cols"`
	Rows  uint16 `json:"rows"`
}

func (options *OpenShellReq) ApplyDefaults() {
	if options.Shell == "" {
		options.Shell = "bash"
	}
	if options.Rows < 10 {
		options.Rows = 10
	}
	if options.Cols < 40 {
		options.Cols = 40
	}
}

func (sf OpenShellReq) WriteTo(w io.Writer) (err error) {
	io.WriteString(w, "OpenShell:")
	e := json.NewEncoder(w)
	err = e.Encode(sf)
	return
}
