package util

import (
	"io"
)

type Rfc2045 struct {
	w     io.Writer
	count int
	write int
}

func NewRfc2045(size int, w io.Writer) *Rfc2045 {
	return &Rfc2045{w: w, count: size, write: size}
}

func (rfc *Rfc2045) Write(data []byte) (total int, err error) {
	var chunk int
	for len(data) > rfc.write {
		if chunk, err = rfc.w.Write(data[:rfc.write]); err != nil {
			return
		}
		data = data[rfc.write:]
		total += chunk
		rfc.write = rfc.count
		rfc.w.Write([]byte{'\n'})
	}
	chunk, err = rfc.w.Write(data)
	total += chunk
	rfc.write -= chunk
	return
}
