package util

import (
	"encoding/base64"
	"io"
	"time"
)

type Rfc2045 struct {
	w io.Writer
	n int
	encoder io.WriteCloser 
}

func NewRfc2045(writer io.Writer) io.WriteCloser {
	self := &Rfc2045{w: writer, n: 76}
	self.encoder = base64.NewEncoder(base64.StdEncoding, self)
	return self.encoder
}

// func (self *Rfc2045) Write(data []byte) (total int, err error) { 
// 	return self.encoder.Write(data)
// }

func (self *Rfc2045) Write(data []byte) (total int, err error) {
	var n int
	for len(data) > 0 {
		if len(data) < self.n {
			if n, err = self.w.Write(data); err != nil {
				break
			}
			total += n
			self.n -= n
			data = []byte{}
		} else {
			if n, err = self.w.Write(data[0:self.n]); err != nil {
				break
			}
			self.w.Write([]byte{'\n'})
			total += n + 1
			data = data[self.n:]
			self.n = 76
		}
		time.Sleep(time.Microsecond)
	}
	return
}
