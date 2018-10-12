package transport

import (
	"errors"
	"net"

	"github.com/nknorg/nnet/config"
)

// Transport is an abstract transport layer between local and remote nodes
type Transport interface {
	Dial(raddr string) (net.Conn, error)
	Listen(port uint16) (net.Listener, error)
	GetProtocol() string
}

// NewTransport creates a transport based on conf
func NewTransport(conf *config.Config) (Transport, error) {
	switch conf.Transport {
	case "kcp":
		return NewKCPTransport(), nil
	case "tcp":
		return NewTCPTransport(conf.DialTimeout), nil
	default:
		return nil, errors.New("Unknown transport" + conf.Transport)
	}
}
