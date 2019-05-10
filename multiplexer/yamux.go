package multiplexer

import (
	"net"

	"github.com/hashicorp/yamux"
)

// Yamux is the multiplexer using yamux
type Yamux struct {
	session *yamux.Session
}

// NewYamux creates a new yamux multiplexer
func NewYamux(conn net.Conn, isClient bool) (Multiplexer, error) {
	var session *yamux.Session
	var err error
	if isClient {
		session, err = yamux.Client(conn, nil)
	} else {
		session, err = yamux.Server(conn, nil)
	}
	if err != nil {
		return nil, err
	}

	return &Yamux{
		session: session,
	}, nil
}

// AcceptStream is used to block until the next available stream is ready to be
// accepted
func (s *Yamux) AcceptStream() (net.Conn, error) {
	return s.session.AcceptStream()
}

// OpenStream is used to create a new stream
func (s *Yamux) OpenStream() (net.Conn, error) {
	return s.session.OpenStream()
}
