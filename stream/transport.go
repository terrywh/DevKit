package stream

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"time"

	quic "github.com/quic-go/quic-go"
	"github.com/terrywh/devkit/entity"
	"github.com/terrywh/devkit/infra/log"
)

func DeviceIDFromCert(cert []byte) entity.DeviceID {
	hash := sha256.New()
	hash.Write(cert)
	return entity.DeviceID(fmt.Sprintf("%02x", hash.Sum(nil)))
}

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
		panic("failed to init transport: " + err.Error())
	}
	if conn, err = net.ListenUDP("udp", addr); err != nil {
		panic("failed to init transport: " + err.Error())
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

func (tr *Transport) createListener(options *ServerOptions) (*quic.Listener, error) {
	cert, err := tls.LoadX509KeyPair(options.Certificate, options.PrivateKey)
	if err != nil {
		return nil, err
	}
	return tr.transport.Listen(&tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{options.ApplicationProtocol},
		ClientAuth:   tls.RequireAnyClientCert,
		VerifyPeerCertificate: func(certs [][]byte, chains [][]*x509.Certificate) error {

			for _, cert := range certs {
				device_id := DeviceIDFromCert(cert)
				if options.Authorize(device_id) {
					return nil
				}
			}

			return entity.ErrUnauthorized
		},
	}, &quic.Config{
		HandshakeIdleTimeout: 10 * time.Second,
		KeepAlivePeriod:      25 * time.Second,
		Allow0RTT:            true,
		EnableDatagrams:      true,
	})
}

type DialOptions struct {
	Address     string
	Certificate string
	PrivateKey  string

	ApplicationProtocol string
	Retry               int           // 默认 0 时，不做重试；当 Retry < 0 时无限重试
	Backoff             time.Duration // 默认 2400ms 重试间隔
}

func (tr *Transport) dial(ctx context.Context, options *DialOptions) (conn quic.Connection, device_id entity.DeviceID, err error) {
	if err = ctx.Err(); err != nil {
		return
	}
	var cert tls.Certificate
	cert, err = tls.LoadX509KeyPair(options.Certificate, options.PrivateKey)
	if err != nil {
		return nil, "", err
	}
	var addr *net.UDPAddr
	addr, err = net.ResolveUDPAddr("udp", options.Address)
	if err != nil {
		return nil, "", err
	}
	conn, err = tr.transport.Dial(ctx, addr, &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"devkit"},
		Certificates:       []tls.Certificate{cert},
		VerifyPeerCertificate: func(certs [][]byte, chain [][]*x509.Certificate) error {
			device_id = DeviceIDFromCert(certs[0])
			return nil
		},
	}, &quic.Config{
		HandshakeIdleTimeout: 10 * time.Second,
		KeepAlivePeriod:      25 * time.Second,
		Allow0RTT:            true,
		EnableDatagrams:      true,
	})
	return
}

func (tr *Transport) Dial(ctx context.Context, options *DialOptions) (conn quic.Connection, device_id entity.DeviceID, err error) {
	for i := 0; i < options.Retry; i++ {
		if err = ctx.Err(); err != nil {
			break
		}
		log.Trace("-> ", options.Address, "(", i+1, "/", options.Retry, ")")
		if conn, device_id, err = tr.dial(ctx, options); err == nil && conn != nil {
			break
		}
		time.Sleep(options.Backoff)
	}
	return
}

func (tr *Transport) Close() (err error) {
	if tr.transport != nil {
		err = tr.transport.Close()
		tr.transport = nil
	}
	return
}

func (tr *Transport) WriteTo(b []byte, a net.Addr) (int, error) {
	return tr.transport.WriteTo(b, a)
}
