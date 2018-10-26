package transport

import (
	"fmt"
	"net"

	kcp "github.com/xtaci/kcp-go"
)

// KCPTransport is the transport layer based on KCP protocol
type KCPTransport struct {
	protocol string
}

// NewKCPTransport creates a new KCP transport layer
func NewKCPTransport() *KCPTransport {
	t := &KCPTransport{
		protocol: "udp",
	}
	return t
}

// Dial connects to the remote address on the network "udp"
func (t *KCPTransport) Dial(addr string) (net.Conn, error) {
	return kcp.Dial(addr)
}

// Listen listens for incoming packets to "port" on the network "udp"
func (t *KCPTransport) Listen(port uint16) (net.Listener, error) {
	laddr := fmt.Sprintf(":%d", port)
	return kcp.Listen(laddr)
}

// GetProtocol returns the network protocol used (tcp or udp)
func (t *KCPTransport) GetProtocol() string {
	return t.protocol
}
