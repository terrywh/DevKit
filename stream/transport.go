package stream

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	quic "github.com/quic-go/quic-go"
)

type TransportOptions struct {
	LocalAddress string
}

type Transport struct {
	transport *quic.Transport
}

var DefaultTransport *Transport

func InitTransport(options TransportOptions) (tr *Transport, err error) {
	var conn *net.UDPConn
	var addr *net.UDPAddr
	if options.LocalAddress == "" {
		options.LocalAddress = "0.0.0.0:0"
	}

	if addr, err = net.ResolveUDPAddr("udp", options.LocalAddress); err != nil {
		return
	}
	if conn, err = net.ListenUDP("udp", addr); err != nil {
		return
	}
	tr = &Transport{
		transport: &quic.Transport{
			Conn: conn,
		},
	}
	DefaultTransport = tr
	return
}

func (tr *Transport) LocalAddress() net.Addr {
	return tr.transport.Conn.LocalAddr()
}

type ServerOptions struct {
	Certificate string
	PrivateKey  string

	ApplicationProtocol string
}

func (so *ServerOptions) ApplyDefaults() {
	if so.Certificate == "" || so.PrivateKey == "" {
		so.Certificate = "./var/cert/server.crt"
		so.PrivateKey = "./var/cert/server.key"
	}
	if so.ApplicationProtocol == "" {
		so.ApplicationProtocol = "devkit"
	}
}

func (tr *Transport) CreateServer(options ServerOptions) (*quic.Listener, error) {
	options.ApplyDefaults()
	cert, err := tls.LoadX509KeyPair(options.Certificate, options.PrivateKey)
	if err != nil {
		return nil, err
	}
	return tr.transport.Listen(&tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{options.ApplicationProtocol},
	}, &quic.Config{
		KeepAlivePeriod: 25 * time.Second,
		Allow0RTT:       true,
	})
}

type DialOptions struct {
	Address string
	Retry   int           // 默认 0 时，不做重试；当 Retry < 0 时无限重试
	Backoff time.Duration // 默认 3s 重试间隔
}

func (options *DialOptions) ApplyDefaults() {
	if options.Backoff < time.Second {
		options.Backoff = 3 * time.Second
	}
}

func (tr *Transport) Dial(ctx context.Context, address string) (quic.Connection, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}
	return tr.transport.Dial(ctx, addr, &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"devkit"},
	}, &quic.Config{
		KeepAlivePeriod: 25 * time.Second,
		Allow0RTT:       true,
	})
}

func (tr *Transport) DialEx(ctx context.Context, options DialOptions) (conn quic.Connection, err error) {
	for i := 0; i < options.Retry; i++ {
		if conn, err = tr.Dial(ctx, options.Address); err == nil && conn != nil {
			break
		}
		time.Sleep(options.Backoff)
	}
	return
}

func (tr *Transport) Close() error {
	return tr.transport.Close()
}
