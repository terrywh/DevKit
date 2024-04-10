package infra

import "io"

type Pseudo interface {
	io.ReadWriteCloser
	Resize(cols, rows int) error
	Pid() int
}
