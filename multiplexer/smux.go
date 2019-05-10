package multiplexer

import (
	"net"

	"github.com/xtaci/smux"
)

// Smux is the multiplexer using smux
type Smux struct {
	session *smux.Session
}

// NewSmux creates a new smux multiplexer
func NewSmux(conn net.Conn, isClient bool) (Multiplexer, error) {
	var session *smux.Session
	var err error
	if isClient {
		session, err = smux.Client(conn, nil)
	} else {
		session, err = smux.Server(conn, nil)
	}
	if err != nil {
		return nil, err
	}

	return &Smux{
		session: session,
	}, nil
}

// AcceptStream is used to block until the next available stream is ready to be
// accepted
func (s *Smux) AcceptStream() (net.Conn, error) {
	return s.session.AcceptStream()
}

// OpenStream is used to create a new stream
func (s *Smux) OpenStream() (net.Conn, error) {
	return s.session.OpenStream()
}
