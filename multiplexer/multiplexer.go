package multiplexer

import (
	"errors"
	"net"
)

// Multiplexer is the multiplexer interface
type Multiplexer interface {
	AcceptStream() (net.Conn, error)
	OpenStream() (net.Conn, error)
}

// NewMultiplexer creates a multiplexer based on conf
func NewMultiplexer(protocol string, conn net.Conn, isClient bool) (Multiplexer, error) {
	switch protocol {
	case "smux":
		return NewSmux(conn, isClient)
	case "yamux":
		return NewYamux(conn, isClient)
	default:
		return nil, errors.New("Unknown protocol " + protocol)
	}
}
