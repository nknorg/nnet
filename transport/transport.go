package transport

import (
	"errors"
	"net"
	"time"
)

// Transport is an abstract transport layer between local and remote nodes
type Transport interface {
	Dial(addr string, dialTimeout time.Duration) (net.Conn, error)
	Listen(port uint16) (net.Listener, error)
	GetNetwork() string
	String() string
}

// NewTransport creates a transport based on conf
func NewTransport(protocol string) (Transport, error) {
	switch protocol {
	case "kcp":
		return NewKCPTransport(), nil
	case "tcp":
		return NewTCPTransport(), nil
	default:
		return nil, errors.New("Unknown protocol " + protocol)
	}
}
