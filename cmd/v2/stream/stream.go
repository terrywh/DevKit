package stream

import (
	"context"
	"crypto/tls"
	"net"
	"sync"
	"time"

	quic "github.com/quic-go/quic-go"
)

var localaddr net.Addr
var transport *quic.Transport
var once sync.Once

func initByServer(options ServerOptions) {
	localaddr, _ = net.ResolveUDPAddr("udp", options.Bind)
	conn, _ := net.ListenUDP("udp", localaddr.(*net.UDPAddr))
	transport = &quic.Transport{
		Conn: conn,
	}
}

func initByClient() {
	addr, _ := net.ResolveUDPAddr("udp", "0.0.0.0:0")
	conn, _ := net.ListenUDP("udp", addr)
	localaddr = conn.LocalAddr()
	transport = &quic.Transport{
		Conn: conn,
	}
}

type ServerOptions struct {
	Bind string
	Crt  string
	Key  string
}

func (so *ServerOptions) MakeDefaults() {
	if so.Bind == "" {
		so.Bind = "0.0.0.0:12345"
	}
	if so.Crt == "" || so.Key == "" {
		so.Crt = "./var/cert/server.crt"
		so.Key = "./var/cert/server.key"
	}
}

func CreateServer(options ServerOptions) (*quic.Listener, error) {
	once.Do(func() {
		initByServer(options)
	})
	options.MakeDefaults()
	cert, err := tls.LoadX509KeyPair(options.Crt, options.Key)
	if err != nil {
		return nil, err
	}
	return transport.Listen(&tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"devkit"},
	}, &quic.Config{
		KeepAlivePeriod: 25 * time.Second,
		Allow0RTT:       true,
	})
}

func Connect(ctx context.Context, addr net.Addr) (quic.Connection, error) {
	once.Do(initByClient)
	// addr, _ := net.ResolveUDPAddr("udp", address)
	return transport.Dial(ctx, addr, &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"devkit"},
	}, &quic.Config{
		KeepAlivePeriod: 25 * time.Second,
	})
}

func LocalAddress() net.Addr {
	return localaddr
}
