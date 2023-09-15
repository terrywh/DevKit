package ssh

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

type ConnectionRequest struct {
	Host string `json:"host"`
	Port int `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

type Request struct {
	Route []ConnectionRequest `json:"route"`
	Command string `json:"init"`
	Rows int `json:"rows"`
	Cols int `json:"cols"`
}

func (req Request) Key() string {
	hash := md5.New()
	for _, route := range req.Route {
		fmt.Fprint(hash, route.User, "@", route.Host, ":", route.Port)
	}
	data := hash.Sum(nil)
	return hex.EncodeToString(data)
}

