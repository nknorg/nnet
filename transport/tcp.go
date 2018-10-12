package transport

import (
	"fmt"
	"net"
	"time"
)

// TCPTransport is the transport layer based on TCP protocol
type TCPTransport struct {
	protocol    string
	dialTimeout time.Duration
}

// NewTCPTransport creates a new TCP transport layer
func NewTCPTransport(dialTimeout time.Duration) *TCPTransport {
	t := &TCPTransport{
		protocol:    "tcp",
		dialTimeout: dialTimeout,
	}
	return t
}

// Dial connects to the remote address "raddr" on the network "tcp"
func (t *TCPTransport) Dial(raddr string) (net.Conn, error) {
	return net.DialTimeout(t.GetProtocol(), raddr, t.dialTimeout)
}

// Listen listens for incoming packets to "port" on the network "tcp"
func (t *TCPTransport) Listen(port uint16) (net.Listener, error) {
	laddr := fmt.Sprintf(":%d", port)
	return net.Listen(t.GetProtocol(), laddr)
}

// GetProtocol returns the network protocol used (tcp or udp)
func (t *TCPTransport) GetProtocol() string {
	return t.protocol
}
