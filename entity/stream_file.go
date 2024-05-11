package entity

import (
	"crypto"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

type File struct {
	Path string `json:"path"`
	Size int64  `json:"size,omitempty"`
	Perm uint32 `json:"perm,omitempty"`
}

type SelectFile struct {
	Path   string `json:"file"`
	Expire int64  `json:"expire"`
	Auth   string `json:"auth"`
}

func (fa SelectFile) GenSign(signer crypto.Signer) string {
	hash := sha256.New()
	fmt.Fprintf(hash, "%s:%d", fa.Path, fa.Expire)

	sign, err := signer.Sign(rand.Reader, hash.Sum(nil), crypto.SHA256)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(sign)
}

type StreamFile struct {
	Source  File              `json:"source,omitempty"`
	Target  File              `json:"target,omitempty"`
	Options StreamFileOptions `json:"options"`
}

type StreamFileOptions struct {
	Override bool `json:"override,omitempty"`
}
