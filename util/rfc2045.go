package util

import (
	"io"
	"time"
)

type Rfc2045 struct {
	w      io.Writer
	c      int
	n      int
	length int
	buffer []byte
}

func NewRfc2045(size int, w io.Writer) *Rfc2045 {
	return &Rfc2045{w: w, c: size, n: size, buffer: make([]byte, size)}
}

func (self *Rfc2045) Write(data []byte) (total int, err error) {
	var n int
	n = self.Flush()
	self.n -= n
	total += n
	for len(data) >= self.n {
		if n, err = self.w.Write(data[0:self.n]); err != nil {
			break
		}
		total += n
		self.w.Write([]byte{'\n'})
		time.Sleep(10 * time.Microsecond)
		data = data[n:]
		self.n = self.c
	}
	self.length = len(data)
	copy(self.buffer, data)
	return
}

func (self *Rfc2045) Flush() (total int) {
	var n int
	if len(self.buffer) > 0 {
		n, _ = self.w.Write(self.buffer[:self.length])
		self.length = 0
		total += n
	}
	return
}
