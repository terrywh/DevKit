package util

import (
	"io"
	"os"
	"testing"
)

func TestRfc2045(t *testing.T) {
	source, _ := os.Open("rfc2045.go")
	defer source.Close()
	target, _ := os.Create("rfc2045.txt")
	defer target.Close()
	encoder := NewRfc2045(target)
	io.Copy(encoder, source)
	encoder.Close()
}