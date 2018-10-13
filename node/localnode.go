package node

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/nknorg/nnet/cache"
	"github.com/nknorg/nnet/config"
	"github.com/nknorg/nnet/log"
	"github.com/nknorg/nnet/protobuf"
	"github.com/nknorg/nnet/transport"
)

const (
	// Max number of msg that can be buffered per routing type
	rxMsgChanLen = 23333

	// Max number of msg to be processed that can be buffered
	handleMsgChanLen = 23333

	// How long a reply chan becomes expired after created
	replyChanExpiration = replyTimeout

	// How often to check and delete expired reply chan
	replyChanCleanupInterval = 1 * time.Second
)

// LocalNode is a local node
type LocalNode struct {
	*Node
	transport      transport.Transport
	port           uint16
	listener       net.Listener
	rxMsgChan      map[protobuf.RoutingType]chan *RemoteMessage
	handleMsgChan  chan *RemoteMessage
	replyChanCache cache.Cache
	*middlewareStore
}

// NewLocalNode creates a local node
func NewLocalNode(id []byte, conf *config.Config) (*LocalNode, error) {
	if id == nil {
		return nil, errors.New("node id is nil")
	}

	trans, err := transport.NewTransport(conf)
	if err != nil {
		return nil, err
	}

	port := conf.Port

	addr := fmt.Sprintf(":%d", port)

	node, err := NewNode(id, addr)
	if err != nil {
		return nil, err
	}

	rxMsgChan := make(map[protobuf.RoutingType]chan *RemoteMessage)
	for routingType := range protobuf.RoutingType_name {
		rxMsgChan[protobuf.RoutingType(routingType)] = make(chan *RemoteMessage, rxMsgChanLen)
	}

	handleMsgChan := make(chan *RemoteMessage, handleMsgChanLen)

	replyChanCache := cache.NewGoCache(replyChanExpiration, replyChanCleanupInterval)

	middlewareStore := newMiddlewareStore()

	localNode := &LocalNode{
		Node:            node,
		transport:       trans,
		port:            port,
		rxMsgChan:       rxMsgChan,
		handleMsgChan:   handleMsgChan,
		replyChanCache:  replyChanCache,
		middlewareStore: middlewareStore,
	}

	return localNode, nil
}

// Start starts the runtime loop of the local node
func (ln *LocalNode) Start() error {
	ln.StartOnce.Do(func() {
		go ln.handleMsg()
		go ln.listen()
	})

	return nil
}

// Stop stops the local node
func (ln *LocalNode) Stop(err error) {
	ln.StopOnce.Do(func() {
		if err != nil {
			log.Warnf("Local node %v stops because of error: %s", ln.Node, err)
		} else {
			log.Infof("Local node %v stops", ln.Node)
		}

		ln.LifeCycle.Stop()

		if ln.listener != nil {
			ln.listener.Close()
		}
	})
}

// handleMsg starts a loop that handles received msg
func (ln *LocalNode) handleMsg() {
	var remoteMsg *RemoteMessage
	var err error

	for {
		if ln.IsStopped() {
			return
		}

		select {
		case remoteMsg = <-ln.handleMsgChan:
			err = ln.handleRemoteMessage(remoteMsg)
			if err != nil {
				log.Error(err)
				continue
			}
		}
	}
}

// listen listens for incoming connections
func (ln *LocalNode) listen() {
	listener, err := ln.transport.Listen(ln.port)
	if err != nil {
		ln.Stop(fmt.Errorf("failed to listen to port %d", ln.port))
		return
	}
	ln.listener = listener

	for {
		// listener.Accept() is placed before checking stops to prevent the error
		// log when local node is stopped and thus conn is closed
		conn, err := listener.Accept()

		if ln.IsStopped() {
			return
		}

		if err != nil {
			log.Error("Error accepting connection:", err)
			time.Sleep(1 * time.Second)
			continue
		}

		log.Infof("Remote node connect from %s to local address %s", conn.RemoteAddr(), conn.LocalAddr())

		_, err = ln.StartRemoteNode(conn, false)
		if err != nil {
			log.Error("Error creating remote node:", err)
			continue
		}
	}
}

// Connect try to establish connection with address remoteNodeAddr
func (ln *LocalNode) Connect(remoteNodeAddr string) error {
	if remoteNodeAddr == ln.Node.Addr {
		return errors.New("trying to connect to self")
	}

	conn, err := ln.transport.Dial(remoteNodeAddr)
	if err != nil {
		return err
	}

	_, err = ln.StartRemoteNode(conn, true)
	if err != nil {
		return err
	}

	return nil
}

// StartRemoteNode creates and starts a remote node using conn
func (ln *LocalNode) StartRemoteNode(conn net.Conn, isOutbound bool) (*RemoteNode, error) {
	remoteNode, err := NewRemoteNode(ln, conn, isOutbound)
	if err != nil {
		return nil, err
	}

	for _, f := range ln.middlewareStore.remoteNodeConnected {
		if !f(remoteNode) {
			break
		}
	}

	err = remoteNode.Start()
	if err != nil {
		return nil, err
	}

	return remoteNode, nil
}

// GetRxMsgChan gets the message channel of a routing type, or return error if
// channel for routing type does not exist
func (ln *LocalNode) GetRxMsgChan(routingType protobuf.RoutingType) (chan *RemoteMessage, error) {
	c, ok := ln.rxMsgChan[routingType]
	if !ok {
		return nil, fmt.Errorf("Msg chan does not exist for type %d", routingType)
	}
	return c, nil
}

// AllocReplyChan creates a reply chan for msg with id msgID
func (ln *LocalNode) AllocReplyChan(msgID []byte) (chan *RemoteMessage, error) {
	if msgID == nil || len(msgID) == 0 {
		return nil, errors.New("Message id is empty")
	}

	replyChan := make(chan *RemoteMessage)

	err := ln.replyChanCache.Add(msgID, replyChan)
	if err != nil {
		return nil, err
	}

	return replyChan, nil
}

// GetReplyChan gets the message reply channel for message id msgID
func (ln *LocalNode) GetReplyChan(msgID []byte) (chan *RemoteMessage, bool) {
	value, ok := ln.replyChanCache.Get(msgID)
	if !ok {
		return nil, false
	}

	replyChan, ok := value.(chan *RemoteMessage)
	if !ok {
		return nil, false
	}

	return replyChan, true
}

// HandleRemoteMessage add remoteMsg to handleMsgChan for further processing
func (ln *LocalNode) HandleRemoteMessage(remoteMsg *RemoteMessage) error {
	select {
	case ln.handleMsgChan <- remoteMsg:
	default:
		log.Warnf("Local node handle msg chan full, discarding msg")
	}
	return nil
}
