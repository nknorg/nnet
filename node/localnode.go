package node

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/nknorg/nnet/cache"
	"github.com/nknorg/nnet/config"
	"github.com/nknorg/nnet/log"
	"github.com/nknorg/nnet/protobuf"
	"github.com/nknorg/nnet/transport"
)

const (
	// How many concurrent goroutines are handling messages
	numWorkers = 1
)

// LocalNode is a local node
type LocalNode struct {
	*Node
	*config.Config
	*middlewareStore
	address        *transport.Address
	port           uint16
	listener       net.Listener
	handleMsgChan  chan *RemoteMessage
	rxMsgChan      map[protobuf.RoutingType]chan *RemoteMessage
	rxMsgCache     cache.Cache
	replyChanCache cache.Cache
	replyTimeout   time.Duration
	neighbors      sync.Map
}

// NewLocalNode creates a local node
func NewLocalNode(id []byte, conf *config.Config) (*LocalNode, error) {
	if id == nil {
		return nil, errors.New("node id is nil")
	}

	address, err := transport.NewAddress(conf.Transport, conf.Hostname, conf.Port)
	if err != nil {
		return nil, err
	}

	node, err := NewNode(id, address.String())
	if err != nil {
		return nil, err
	}

	handleMsgChan := make(chan *RemoteMessage, conf.LocalHandleMsgChanLen)

	rxMsgChan := make(map[protobuf.RoutingType]chan *RemoteMessage)

	rxMsgCache := cache.NewGoCache(conf.LocalRxMsgCacheExpiration, conf.LocalRxMsgCacheCleanupInterval)

	replyChanCache := cache.NewGoCache(conf.DefaultReplyTimeout, conf.ReplyChanCleanupInterval)

	middlewareStore := newMiddlewareStore()

	localNode := &LocalNode{
		Node:            node,
		Config:          conf,
		middlewareStore: middlewareStore,
		address:         address,
		port:            conf.Port,
		handleMsgChan:   handleMsgChan,
		rxMsgChan:       rxMsgChan,
		rxMsgCache:      rxMsgCache,
		replyChanCache:  replyChanCache,
		replyTimeout:    conf.DefaultReplyTimeout,
	}

	for routingType := range protobuf.RoutingType_name {
		localNode.RegisterRoutingType(protobuf.RoutingType(routingType))
	}

	return localNode, nil
}

// Start starts the runtime loop of the local node
func (ln *LocalNode) Start() error {
	ln.StartOnce.Do(func() {
		for _, mw := range ln.middlewareStore.localNodeWillStart {
			if !mw.Func(ln) {
				break
			}
		}

		for i := 0; i < numWorkers; i++ {
			go ln.handleMsg()
		}

		go ln.listen()

		for _, mw := range ln.middlewareStore.localNodeStarted {
			if !mw.Func(ln) {
				break
			}
		}
	})

	return nil
}

// Stop stops the local node
func (ln *LocalNode) Stop(err error) {
	ln.StopOnce.Do(func() {
		for _, mw := range ln.middlewareStore.localNodeWillStop {
			if !mw.Func(ln) {
				break
			}
		}

		if err != nil {
			log.Warningf("Local node %v stops because of error: %s", ln, err)
		} else {
			log.Infof("Local node %v stops", ln)
		}

		ln.neighbors.Range(func(key, value interface{}) bool {
			remoteNode, ok := value.(*RemoteNode)
			if ok {
				remoteNode.Stop(err)
			}
			return true
		})

		time.Sleep(stopGracePeriod)

		ln.LifeCycle.Stop()

		if ln.listener != nil {
			ln.listener.Close()
		}

		for _, mw := range ln.middlewareStore.localNodeStopped {
			if !mw.Func(ln) {
				break
			}
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

		remoteMsg = <-ln.handleMsgChan

		err = ln.handleRemoteMessage(remoteMsg)
		if err != nil {
			log.Errorf("Handle remote message error: %v", err)
			continue
		}

	}
}

// listen listens for incoming connections
func (ln *LocalNode) listen() {
	listener, err := ln.address.Transport.Listen(ln.port)
	if err != nil {
		ln.Stop(fmt.Errorf("failed to listen to port %d", ln.port))
		return
	}
	ln.listener = listener

	if ln.port == 0 {
		_, portStr, err := net.SplitHostPort(listener.Addr().String())
		if err != nil {
			ln.Stop(err)
		}

		port, err := strconv.Atoi(portStr)
		if err != nil {
			ln.Stop(err)
		}

		ln.SetInternalPort(uint16(port))
	}

	if ln.address.Port == 0 {
		ln.address.Port = ln.port
		ln.Addr = ln.address.String()
	}

	for {
		// listener.Accept() is placed before checking stops to prevent the error
		// log when local node is stopped and thus conn is closed
		conn, err := listener.Accept()

		if ln.IsStopped() {
			if conn != nil {
				conn.Close()
			}
			return
		}

		if err != nil {
			log.Errorf("Error accepting connection: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		_, loaded := ln.neighbors.LoadOrStore(conn.RemoteAddr().String(), nil)
		if loaded {
			log.Errorf("Remote addr %s is already connected, reject connection", conn.RemoteAddr().String())
			conn.Close()
			continue
		}

		log.Infof("Remote node connect from %s to local address %s", conn.RemoteAddr().String(), conn.LocalAddr())

		rn, err := ln.StartRemoteNode(conn, false, nil)
		if err != nil {
			log.Errorf("Error creating remote node: %v", err)
			ln.neighbors.Delete(conn.RemoteAddr().String())
			conn.Close()
			continue
		}

		ln.neighbors.Store(conn.RemoteAddr().String(), rn)
	}
}

// SetInternalPort changes the internal port that local node will listen on. It should
// not be called once the node starts. Note that it does not change the external
// port that other nodes will connect to.
func (ln *LocalNode) SetInternalPort(port uint16) {
	ln.port = port
}

// Connect try to establish connection with address remoteNodeAddr, returns the
// remote node, if the remote node is ready, and error. The remote rode can be
// nil if another goroutine is connecting to the same address concurrently. The
// remote node is ready if an active connection to the remoteNodeAddr exists and
// node info has been exchanged.
func (ln *LocalNode) Connect(n *protobuf.Node) (*RemoteNode, bool, error) {
	if n.Addr == ln.address.String() {
		return nil, false, errors.New("trying to connect to self")
	}

	remoteAddress, err := transport.Parse(n.Addr)
	if err != nil {
		return nil, false, err
	}

	key := remoteAddress.ConnRemoteAddr()
	value, loaded := ln.neighbors.LoadOrStore(key, nil)
	if loaded {
		remoteNode, ok := value.(*RemoteNode)
		if ok {
			if remoteNode.IsStopped() {
				log.Warningf("Remove stopped remote node %v from list", remoteNode)
				ln.neighbors.Delete(key)
			} else {
				log.Infof("Load remote node %v from list", remoteNode)
				return remoteNode, remoteNode.IsReady(), nil
			}
		} else {
			log.Infof("Another goroutine is connecting to %s", key)
			return nil, false, nil
		}
	}

	var shouldConnect, shouldCallNextMiddleware bool
	for _, mw := range ln.middlewareStore.willConnectToNode {
		shouldConnect, shouldCallNextMiddleware = mw.Func(n)
		if !shouldConnect {
			return nil, false, nil
		}
		if !shouldCallNextMiddleware {
			break
		}
	}

	conn, err := remoteAddress.Dial(ln.DialTimeout)
	if err != nil {
		ln.neighbors.Delete(key)
		return nil, false, err
	}

	remoteNode, err := ln.StartRemoteNode(conn, true, n)
	if err != nil {
		ln.neighbors.Delete(key)
		conn.Close()
		return nil, false, err
	}

	ln.neighbors.Store(key, remoteNode)

	return remoteNode, false, nil
}

// StartRemoteNode creates and starts a remote node using conn. If n is not
// nil, its id and address will be used to validate ExchangeNode info.
func (ln *LocalNode) StartRemoteNode(conn net.Conn, isOutbound bool, n *protobuf.Node) (*RemoteNode, error) {
	remoteNode, err := NewRemoteNode(ln, conn, isOutbound, n)
	if err != nil {
		return nil, err
	}

	for _, mw := range ln.middlewareStore.remoteNodeConnected {
		if !mw.Func(remoteNode) {
			break
		}
	}

	err = remoteNode.Start()
	if err != nil {
		return nil, err
	}

	return remoteNode, nil
}

// RegisterRoutingType register a routing type and creates the rxMsgChan for it
func (ln *LocalNode) RegisterRoutingType(routingType protobuf.RoutingType) {
	ln.rxMsgChan[routingType] = make(chan *RemoteMessage, ln.LocalRxMsgChanLen)
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
func (ln *LocalNode) AllocReplyChan(msgID []byte, expiration time.Duration) (chan *RemoteMessage, error) {
	if len(msgID) == 0 {
		return nil, errors.New("Message id is empty")
	}

	replyChan := make(chan *RemoteMessage)

	err := ln.replyChanCache.AddWithExpiration(msgID, replyChan, expiration)
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

// AddToRxCache add RemoteMessage id to rxMsgCache if not exists. Returns if msg
// id is added (instead of loaded) and error when adding
func (ln *LocalNode) AddToRxCache(msgID []byte) (bool, error) {
	_, found := ln.rxMsgCache.Get(msgID)
	if found {
		return false, nil
	}

	err := ln.rxMsgCache.Add(msgID, struct{}{})
	if err != nil {
		if _, found := ln.rxMsgCache.Get(msgID); found {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// HandleRemoteMessage add remoteMsg to handleMsgChan for further processing
func (ln *LocalNode) HandleRemoteMessage(remoteMsg *RemoteMessage) error {
	select {
	case ln.handleMsgChan <- remoteMsg:
	default:
		log.Warningf("Local node handle msg chan full, discarding msg")
	}
	return nil
}

// GetNeighbors returns a list of remote nodes that are connected to local nodes
// where the filter function returns true. Pass nil filter to return all
// neighbors.
func (ln *LocalNode) GetNeighbors(filter func(*RemoteNode) bool) ([]*RemoteNode, error) {
	nodes := make([]*RemoteNode, 0)
	ln.neighbors.Range(func(key, value interface{}) bool {
		remoteNode, ok := value.(*RemoteNode)
		if ok && remoteNode.IsReady() && !remoteNode.IsStopped() {
			if filter == nil || filter(remoteNode) {
				nodes = append(nodes, remoteNode)
			}
		}
		return true
	})
	return nodes, nil
}
