package util

import (
	"encoding/base64"
	"io"
	"os"
	"testing"
)

func TestRfc2045(t *testing.T) {
	source, _ := os.Open("../var/trzsz_1.1.7_linux_x86_64.tar.gz")
	defer source.Close()

	target, _ := os.Create("rfc2045.txt")
	defer target.Close()

	rfc2045 := NewRfc2045(128, target)
	encoder := base64.NewEncoder(base64.StdEncoding, rfc2045)

	io.Copy(encoder, source)
	encoder.Close()
}
