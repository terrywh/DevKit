package stream

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
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

type ServerHandler interface {
	ServeStream(ctx context.Context, s quic.Stream)
}

type ServerOptions struct {
	Handler     ServerHandler
	Authorize   func(device_id string) bool
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

var ErrUnauthorized error = errors.New("unauthorized")

func (tr *Transport) CreateServer(options ServerOptions) (*Server, error) {
	options.ApplyDefaults()
	cert, err := tls.LoadX509KeyPair(options.Certificate, options.PrivateKey)
	if err != nil {
		return nil, err
	}
	l, err := tr.transport.Listen(&tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{options.ApplicationProtocol},
		ClientAuth:   tls.RequireAnyClientCert,
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {

			for _, cert := range rawCerts {
				hash := sha256.New()
				hash.Write(cert)
				if options.Authorize(fmt.Sprintf("%2x", hash.Sum(nil))) {
					return nil
				}
			}

			return ErrUnauthorized
		},
	}, &quic.Config{
		KeepAlivePeriod: 25 * time.Second,
		Allow0RTT:       true,
		EnableDatagrams: true,
	})
	if err != nil {
		return nil, err
	}
	return newServer(l, options.Handler)
}

type DialOptions struct {
	Address     string
	Certificate string
	PrivateKey  string

	Retry   int           // 默认 0 时，不做重试；当 Retry < 0 时无限重试
	Backoff time.Duration // 默认 3s 重试间隔
}

func (do *DialOptions) ApplyDefaults() {
	if do.Certificate == "" || do.PrivateKey == "" {
		do.Certificate = "./var/cert/server.crt"
		do.PrivateKey = "./var/cert/server.key"
	}
	if do.Backoff < time.Second {
		do.Backoff = 3 * time.Second
	}
}

func (tr *Transport) dial(ctx context.Context, options DialOptions) (quic.Connection, error) {
	cert, err := tls.LoadX509KeyPair(options.Certificate, options.PrivateKey)
	if err != nil {
		return nil, err
	}
	addr, err := net.ResolveUDPAddr("udp", options.Address)
	if err != nil {
		return nil, err
	}
	return tr.transport.Dial(ctx, addr, &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"devkit"},
		Certificates:       []tls.Certificate{cert},
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {

			for i, cert := range rawCerts {
				hash := sha256.New()
				hash.Write(cert)
				log.Printf("<Transport.Dial> certificate: %d %2x", i, hash.Sum(nil))
			}

			return nil
		},
	}, &quic.Config{
		KeepAlivePeriod: 25 * time.Second,
		Allow0RTT:       true,
		EnableDatagrams: true,
	})
}

func (tr *Transport) Dial(ctx context.Context, options DialOptions) (conn quic.Connection, err error) {
	options.ApplyDefaults()
	for i := 0; i < options.Retry; i++ {
		if conn, err = tr.dial(ctx, options); err == nil && conn != nil {
			break
		}
		time.Sleep(options.Backoff)
	}
	return
}

func (tr *Transport) Close() (err error) {
	if tr.transport != nil {
		err = tr.transport.Close()
	}
	return
}
