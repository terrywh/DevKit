package entity

import (
	"crypto"
	"io"
	"testing"
	"time"
)

type DummySigner struct{}

func (ds DummySigner) Public() crypto.PublicKey {
	return nil
}

func (ds DummySigner) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	return digest, nil
}

func TestSelectFile(t *testing.T) {
	fa := SelectFile{
		Path:   "/tmp/abc",
		Expire: time.Now().Add(10 * time.Second).Unix(),
	}
	ds := DummySigner{}
	fa.Auth = fa.GenSign(ds)
	t.Log(fa)
	if fa.GenSign(ds) != fa.Auth {
		t.Fatal("failed to generate same signature")
	}
}
