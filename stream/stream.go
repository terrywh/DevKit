package stream

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	quic "github.com/quic-go/quic-go"
)

var transport *quic.Transport = &quic.Transport{
	Conn: CreateConn(),
}

func CreateConn() (conn *net.UDPConn) {
	addr, _ := net.ResolveUDPAddr("udp", "0.0.0.0:12345")
	conn, _ = net.ListenUDP("udp", addr)
	return
}

func CreateServer() (*quic.Listener, error) {
	cert, _ := tls.LoadX509KeyPair("./var/cert/server.crt", "./var/cert/server.key")
	return transport.Listen(&tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"devkit"},
	}, &quic.Config{
		KeepAlivePeriod: 25 * time.Second,
		Allow0RTT:       true,
	})
}

func CreateClient() (quic.Connection, error) {
	addr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:12345")
	return transport.Dial(context.Background(), addr, &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"devkit"},
	}, &quic.Config{
		KeepAlivePeriod: 25 * time.Second,
	})
}
