package transport

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Address is a URI for a node
type Address struct {
	Transport Transport
	Host      string
	Port      uint16
}

// NewAddress creates an Address struct with given protocol and address
func NewAddress(protocol, host string, port uint16, supportedTransports []Transport) (*Address, error) {
	transport, err := NewTransport(protocol, supportedTransports)
	if err != nil {
		return nil, err
	}

	addr := &Address{
		Transport: transport,
		Host:      host,
		Port:      port,
	}

	return addr, nil
}

func (addr *Address) String() string {
	return fmt.Sprintf("%s://%s:%d", addr.Transport, addr.Host, addr.Port)
}

// ConnRemoteAddr returns the remote address string that transport can dial
func (addr *Address) ConnRemoteAddr() string {
	return fmt.Sprintf("%s:%d", addr.Host, addr.Port)
}

// Dial dials the remote address using local transport
func (addr *Address) Dial(dialTimeout time.Duration) (net.Conn, error) {
	return addr.Transport.Dial(addr.ConnRemoteAddr(), dialTimeout)
}

// Parse parses a raw addr string into an Address struct
func Parse(rawAddr string, supportedTransports []Transport) (*Address, error) {
	u, err := url.Parse(rawAddr)
	if err != nil {
		return nil, err
	}

	transport, err := NewTransport(u.Scheme, supportedTransports)
	if err != nil {
		return nil, err
	}

	host, portStr, err := SplitHostPort(u.Host)
	if err != nil {
		return nil, err
	}

	port := 0
	if len(portStr) > 0 {
		port, err = strconv.Atoi(portStr)
		if err != nil {
			return nil, err
		}
	}

	addr := &Address{
		Transport: transport,
		Host:      host,
		Port:      uint16(port),
	}

	return addr, nil
}

// SplitHostPort is the same as net.SplitHostPort, except that it allows
// hostport to be host without port.
func SplitHostPort(hostport string) (host, port string, err error) {
	if !strings.Contains(hostport, ":") {
		hostport += ":"
	}
	return net.SplitHostPort(hostport)
}
