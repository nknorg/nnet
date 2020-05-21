package transport

import (
	"fmt"
	"net"
	"time"

	kcp "github.com/xtaci/kcp-go"
)

type kcpListener struct {
	*kcp.Listener
}

func (listener *kcpListener) Accept() (net.Conn, error) {
	conn, err := listener.Listener.AcceptKCP()
	if err != nil {
		return nil, err
	}

	setKCPConn(conn)

	return conn, nil
}

// KCPTransport is the transport layer based on KCP protocol
type KCPTransport struct{}

// NewKCPTransport creates a new KCP transport layer
func NewKCPTransport() *KCPTransport {
	t := &KCPTransport{}
	return t
}

// Dial connects to the remote address on the network "udp"
func (t *KCPTransport) Dial(addr string, dialTimeout time.Duration) (net.Conn, error) {
	conn, err := kcp.DialWithOptions(addr, nil, 0, 0)
	if err != nil {
		return nil, err
	}

	setKCPConn(conn)

	return conn, nil
}

// Listen listens for incoming packets to "port" on the network "udp"
func (t *KCPTransport) Listen(port uint16) (net.Listener, error) {
	laddr := fmt.Sprintf(":%d", port)
	listener, err := kcp.ListenWithOptions(laddr, nil, 0, 0)
	if err != nil {
		return nil, err
	}

	return &kcpListener{listener}, nil
}

// GetNetwork returns the network used (tcp or udp)
func (t *KCPTransport) GetNetwork() string {
	return "udp"
}

func (t *KCPTransport) String() string {
	return "kcp"
}

func setKCPConn(conn *kcp.UDPSession) {
	conn.SetStreamMode(true)
	conn.SetACKNoDelay(true)
	conn.SetNoDelay(1, 20, 2, 0)
}
